package tbox

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

// engine table to struct data define.
type engine struct {
	dsn string  // mysql dsn
	db  *sql.DB // sql DB

	packageName string // convert to struct package name
	pkgPath     string // gen code to dir path
	isOutputCmd bool   // whether to output to cmd,otherwise output to file
	libPath     string // point the library path relate root path

	addTag bool
	tagKey string // the key value of the tag field, the default is `db:xxx`

	// While the first letter of the field is capitalized
	// whether to convert other letters to lowercase, the default false is not converted
	ucFirstOnly bool

	// Table name returned by TaleName func
	// such as xorm, gorm library uses TableName method to return table name
	enableTableNameFunc bool

	// Whether to add json tag, not added by default
	enableJsonTag bool

	noNullField bool
}

// New create an entry for engine
func New(dsn string, opts ...Option) *engine {
	if dsn == "" {
		log.Fatalln("dsn is empty")
	}

	t := &engine{
		dsn:         dsn,
		packageName: "models",
		pkgPath:     "./app/models",
		tagKey:      "db",
		addTag:      true,
		libPath:     "./app/library",
	}

	for _, o := range opts {
		o(t)
	}

	var err error
	t.db, err = t.connect()
	if err != nil {
		log.Fatalln("dsn connect mysql error: ", err)
	}

	if t.pkgPath == "" {
		log.Fatalln("pkg path is empty")
	}

	if !checkPathExist(t.pkgPath) {
		err = os.MkdirAll(t.pkgPath, 0755)
		if err != nil {
			log.Fatalln("create pkg dir error: ", err)
		}
	}

	return t
}

var noNullList = []string{
	"int", "int64", "float", "float32", "json", "string",
}

func (t *engine) getTableCode(tab string, record []columnEntry) []string {
	var header strings.Builder
	header.WriteString(fmt.Sprintf("// Package %s of db entity\n", t.packageName))
	header.WriteString(fmt.Sprintf("// Code generated by go-god/tbox. DO NOT EDIT!!!\n\n"))
	header.WriteString(fmt.Sprintf("package %s\n\n", "models"))

	var importBuf strings.Builder
	var tabInfoBuf strings.Builder
	tabPrefix := t.camelCase(tab)
	tabName := tabPrefix + "Table"
	structName := tabPrefix
	tabInfoBuf.WriteString(fmt.Sprintf("// %s for %s\n", tabName, tab))
	tabInfoBuf.WriteString(fmt.Sprintf("const %s = \"%s\"\n\n", tabName, tab))
	tabInfoBuf.WriteString(fmt.Sprintf("// %s for %s table entity struct.\n", structName, tab))
	tabInfoBuf.WriteString(fmt.Sprintf("type %s struct {\n", structName+"Model"))
	for _, val := range record {
		dataType := getType(val.DataType) // column type
		if dataType == "time.Time" && importBuf.Len() == 0 {
			importBuf.WriteString("import (\n\t\"time\"\n")
			importBuf.WriteString("\t\"gorm.io/driver/mysql\"\n")
			importBuf.WriteString("\t\"gorm.io/gorm/logger\"\n")
			importBuf.WriteString("\t\"gorm.io/gorm\"\n")
			importBuf.WriteString(")\n\n")
		}

		if t.noNullField && strInSet(dataType, noNullList) {
			val.IsNullable = "NO"
		}

		if val.IsNullable == "YES" {
			tabInfoBuf.WriteString(fmt.Sprintf("\t%s\t*%s", t.camelCase(val.Field), dataType))
		} else {
			tabInfoBuf.WriteString(fmt.Sprintf("\t%s\t%s", t.camelCase(val.Field), dataType))
		}

		if t.addTag {
			tabInfoBuf.WriteString("\t")
			tabInfoBuf.WriteString("`")
			if t.enableJsonTag {
				tabInfoBuf.WriteString(fmt.Sprintf("json:\"%s\"", strings.ToLower(val.Field)))
				tabInfoBuf.WriteString(" ")
			}

			if t.tagKey != "" {
				if t.tagKey == "xorm" {
					if val.FieldKey == "PRI" && val.Extra == "auto_increment" {
						tabInfoBuf.WriteString(fmt.Sprintf("%s:\"%s pk autoincr\"", t.tagKey,
							strings.ToLower(val.Field)))
					} else {
						tabInfoBuf.WriteString(fmt.Sprintf("%s:\"%s\"", t.tagKey, strings.ToLower(val.Field)))
					}
				} else if t.tagKey == "gorm" {
					if val.FieldKey == "PRI" && val.Extra == "auto_increment" {
						tabInfoBuf.WriteString(fmt.Sprintf("%s:\"%s primaryKey\"", t.tagKey,
							strings.ToLower(val.Field)))
					} else {
						tabInfoBuf.WriteString(fmt.Sprintf("%s:\"%s\"", t.tagKey, strings.ToLower(val.Field)))
					}
				} else {
					tabInfoBuf.WriteString(fmt.Sprintf("%s:\"%s\"", t.tagKey, strings.ToLower(val.Field)))
				}
			}

			tabInfoBuf.WriteString("`")
		}

		tabInfoBuf.WriteString("\n")
	}

	tabInfoBuf.WriteString("\tConnect          *gorm.DB            `json:\"-\" gorm:\"-\"`\n")
	tabInfoBuf.WriteString("}\n")

	if t.enableTableNameFunc {
		tabInfoBuf.WriteString(fmt.Sprintf("\n// %s for %s", "TableName", tab))
		tabInfoBuf.WriteString(fmt.Sprintf("\nfunc (%s) TableName() string {\n", structName))
		tabInfoBuf.WriteString(fmt.Sprintf("\treturn %s\n", tabName))
		tabInfoBuf.WriteString("}\n")
	}

	// 这里可以改造，很多项目会提前链接db，直接从配置读取
	tabInfoBuf.WriteString(fmt.Sprintf("func New%sModel() *%sModel {\n", t.camelCase(tab), t.camelCase(tab)))
	tabInfoBuf.WriteString(fmt.Sprintf("\t// 这里的链接从连接池里面获取\n"))
	tabInfoBuf.WriteString(fmt.Sprintf("\tdsn := \"root:123456@tcp(127.0.0.1:3306)/app_api?charset=utf8mb4&parseTime=True&loc=Local\"\n"))
	tabInfoBuf.WriteString(fmt.Sprintf("\tdb, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})\n"))
	tabInfoBuf.WriteString(fmt.Sprintf("\tif err != nil {\n"))
	tabInfoBuf.WriteString(fmt.Sprintf("\t\tpanic(err)\n"))
	tabInfoBuf.WriteString(fmt.Sprintf("}\n"))
	tabInfoBuf.WriteString(fmt.Sprintf("\tins := &%sModel{}\n", t.camelCase(tab)))
	tabInfoBuf.WriteString(fmt.Sprintf("\tins.SetConnect(db)\n"))
	tabInfoBuf.WriteString(fmt.Sprintf("\treturn ins\n"))
	tabInfoBuf.WriteString(fmt.Sprintf("}\n\n"))

	tabInfoBuf.WriteString(fmt.Sprintf("func (t *%sModel) GetTableName() string {\n", t.camelCase(tab)))
	tabInfoBuf.WriteString(fmt.Sprintf("\treturn t.TableName()\n"))
	tabInfoBuf.WriteString(fmt.Sprintf("}\n\n"))

	tabInfoBuf.WriteString(fmt.Sprintf("func (t *%sModel) GetId() int64 {\n", t.camelCase(tab)))
	tabInfoBuf.WriteString(fmt.Sprintf("\treturn t.ID\n"))
	tabInfoBuf.WriteString(fmt.Sprintf("}\n\n"))

	tabInfoBuf.WriteString(fmt.Sprintf("func (t *%sModel) GetConnectName() string {\n", t.camelCase(tab)))
	tabInfoBuf.WriteString(fmt.Sprintf("\treturn \"default\"\n"))
	tabInfoBuf.WriteString(fmt.Sprintf("}\n\n"))

	tabInfoBuf.WriteString(fmt.Sprintf("func (t *%sModel) GetConnect() *gorm.DB {\n", t.camelCase(tab)))
	tabInfoBuf.WriteString(fmt.Sprintf("\tif t.Connect == nil {\n"))
	tabInfoBuf.WriteString(fmt.Sprintf("\t\tpanic(\"table of \" + t.GetTableName() + \" connect is nil!\")\n"))
	tabInfoBuf.WriteString(fmt.Sprintf("\t}\n\n"))
	tabInfoBuf.WriteString(fmt.Sprintf("\t// 指定表名\n"))
	tabInfoBuf.WriteString(fmt.Sprintf("\treturn t.Connect.Table(t.GetTableName()).Debug()\n"))
	tabInfoBuf.WriteString(fmt.Sprintf("}\n\n"))

	tabInfoBuf.WriteString(fmt.Sprintf("func (t *%sModel) SetConnect(connect *gorm.DB) {\n", t.camelCase(tab)))
	tabInfoBuf.WriteString(fmt.Sprintf("\tt.Connect = connect\n"))
	tabInfoBuf.WriteString(fmt.Sprintf("}\n\n"))

	tabInfoBuf.WriteString(fmt.Sprintf("func (t *%sModel) GetPrimaryKey() string {\n", t.camelCase(tab)))
	tabInfoBuf.WriteString(fmt.Sprintf("\treturn \"id\"\n"))
	tabInfoBuf.WriteString(fmt.Sprintf("}\n\n"))

	tabInfoBuf.WriteString(fmt.Sprintf("func (t *%sModel) TableName() string {\n", t.camelCase(tab)))
	tabInfoBuf.WriteString(fmt.Sprintf("\treturn \"%s\"\n", tab))
	tabInfoBuf.WriteString(fmt.Sprintf("}\n\n"))

	return []string{
		header.String(), importBuf.String(), tabInfoBuf.String(),
	}
}

// Run exec table to struct
func (t *engine) Run(table ...string) error {
	tabColumns, err := t.GetColumns(table...)
	if err != nil {
		return err
	}

	for tab, record := range tabColumns {
		str := strings.Join(t.getTableCode(tab, record), "")
		if t.isOutputCmd {
			fmt.Println(str)
		}

		if !checkPathExist(t.libPath) {
			dir_path := filepath.Join(t.libPath, t.ucFirst(tab)+"Dao")
			fmt.Println(dir_path)
			//os.Mkdir(dir_path, os.ModePerm)

			err = os.MkdirAll(dir_path, 0755)
			if err != nil {
				log.Fatalln("create pkg dir error: ", err)
			}
		}

		if !checkPathExist(t.pkgPath) {
			fmt.Println(t.pkgPath)
			//os.Mkdir(dir_path, os.ModePerm)

			err = os.MkdirAll(t.pkgPath, 0755)
			if err != nil {
				log.Fatalln("create pkg dir error: ", err)
			}
		}

		// write to file
		var fileName string
		fileName, err = filepath.Abs(filepath.Join(t.pkgPath, t.ucFirst(tab)+"Model.go"))
		if err != nil {
			log.Println("file path error: ", err)
			continue
		}

		err = os.WriteFile(fileName, []byte(str), 0666)
		if err != nil {
			log.Println("gen code error: ", err)
			return err
		}

		t.createLibrary(tab, record)
	}

	return nil
}

type columnEntry struct {
	TableName    string
	Field        string
	DataType     string
	FieldDesc    string
	FieldKey     string
	OrderBy      int
	IsNullable   string
	MaxLength    sql.NullInt64
	NumericPrec  sql.NullInt64
	NumericScale sql.NullInt64
	Extra        string
	FieldComment string
}

var fields = []string{
	"TABLE_NAME as table_name",
	"COLUMN_NAME as field",
	"DATA_TYPE as data_type",
	"COLUMN_TYPE as field_desc",
	"COLUMN_KEY as field_key",
	"ORDINAL_POSITION as order_by",
	"IS_NULLABLE as is_nullable",
	"CHARACTER_MAXIMUM_LENGTH as max_length",
	"NUMERIC_PRECISION as numeric_prec",
	"NUMERIC_SCALE as numeric_scale",
	"EXTRA as extra",
	"COLUMN_COMMENT as field_comment",
}

// GetColumns Get the information fields related to the table according to the table name
func (t *engine) GetColumns(table ...string) (map[string][]columnEntry, error) {
	var sqlStr = "select " + strings.Join(fields, ",") + " from information_schema.COLUMNS " +
		"WHERE table_schema = DATABASE()"
	if len(table) > 0 {
		sqlStr += " AND TABLE_NAME in (" + t.implodeTable(table...) + ")"
	}

	sqlStr += " order by TABLE_NAME asc, ORDINAL_POSITION asc"
	rows, err := t.db.Query(sqlStr)
	if err != nil {
		log.Println("read table information error: ", err.Error())
		return nil, err
	}

	defer rows.Close()

	records := make(map[string][]columnEntry, 20)
	for rows.Next() {
		col := columnEntry{}
		err = rows.Scan(&col.TableName, &col.Field, &col.DataType, &col.FieldDesc,
			&col.FieldKey, &col.OrderBy, &col.IsNullable, &col.MaxLength, &col.NumericPrec, &col.NumericScale,
			&col.Extra, &col.FieldComment,
		)

		if err != nil {
			log.Println("scan column error: ", err)
			continue
		}

		if _, ok := records[col.TableName]; !ok {
			records[col.TableName] = make([]columnEntry, 0, 20)
		}

		records[col.TableName] = append(records[col.TableName], col)
	}

	return records, nil
}

func (t *engine) connect() (*sql.DB, error) {
	db, err := sql.Open("mysql", t.dsn)
	return db, err
}

func (t *engine) implodeTable(table ...string) string {
	var arr []string
	for _, tab := range table {
		arr = append(arr, fmt.Sprintf("'%s'", tab))
	}

	return strings.Join(arr, ",")
}

func (t *engine) camelCase(str string) string {
	var buf strings.Builder
	arr := strings.Split(str, "_")
	if len(arr) > 0 {
		for _, s := range arr {
			if t.ucFirstOnly {
				buf.WriteString(strings.ToUpper(s[0:1]))
				buf.WriteString(strings.ToLower(s[1:]))
			} else {
				lintStr := lintName(strings.Join([]string{strings.ToUpper(s[0:1]), s[1:]}, ""))
				buf.WriteString(lintStr)
			}
		}
	}

	return buf.String()
}

func (t *engine) ucFirst(str string) string {
	str = t.camelCase(str)
	if len(str) == 0 {
		return str
	}

	return strings.ToLower(str[0:1]) + str[1:]
}

func (t *engine) createItemFile(tab string, record []columnEntry) {
	str := strings.Builder{}
	file_name := t.ucFirst(tab) + ".go"
	os.MkdirAll(fmt.Sprintf("%s/%s", t.libPath, t.ucFirst(tab)+"Dao"), 0777)
	file_path := filepath.Join(t.libPath, t.ucFirst(tab)+"Dao", file_name)

	str.WriteString(fmt.Sprintf("package %s\n\n", t.ucFirst(tab)+"Dao"))
	str.WriteString(fmt.Sprintf("import (\n"))
	str.WriteString(fmt.Sprintf("\t\"github.com/gollerxiong/gbox/components\"\n"))
	str.WriteString(fmt.Sprintf("\t\"%s%s\"\n", t.packageName, strings.Trim(t.pkgPath, ".")))
	str.WriteString(fmt.Sprintf(")\n\n"))

	t.createItemStruct(tab, &str)

	fmt.Println(file_path)

	err := os.WriteFile(file_path, []byte(str.String()), 0666)
	if err != nil {
		log.Println("gen code error: ", err)
	}
}

func (t *engine) createItemStruct(tab string, b *strings.Builder) *strings.Builder {
	b.WriteString(fmt.Sprintf("type %s struct {\n", t.camelCase(tab)))
	b.WriteString(fmt.Sprintf("\tcomponents.BaseObject\n"))
	b.WriteString(fmt.Sprintf("}\n\n"))

	t.createItemGetModel(tab, b)
	return b
}

func (t *engine) createItemGetModel(tab string, b *strings.Builder) *strings.Builder {
	b.WriteString(fmt.Sprintf("func (o *%s) GetModel() *models.%s {\n", t.camelCase(tab), t.camelCase(tab)+"Model"))
	b.WriteString(fmt.Sprintf("\treturn o.BaseObject.GetModel().(*models.%s)\n", t.camelCase(tab)+"Model"))
	b.WriteString(fmt.Sprintf("}\n\n"))
	t.createItemLoadById(tab, b)
	return b
}

func (t *engine) createItemLoadById(tab string, b *strings.Builder) *strings.Builder {
	b.WriteString(fmt.Sprintf("func (o *%s) LoadById(id int64) *%s {\n", t.camelCase(tab), t.camelCase(tab)))
	b.WriteString(fmt.Sprintf("\to.BaseObject.LoadById(id)\n"))
	b.WriteString(fmt.Sprintf("\treturn o\n}\n\n"))
	t.createItemNewItem(tab, b)
	return b
}

func (t *engine) createItemNewItem(tab string, b *strings.Builder) *strings.Builder {
	b.WriteString(fmt.Sprintf("func New%s() *%s {\n", t.camelCase(tab), t.camelCase(tab)))
	b.WriteString(fmt.Sprintf("\tins := &%s{}\n", t.camelCase(tab)))
	b.WriteString(fmt.Sprintf("\tins.SetField(\"*\")\n"))
	b.WriteString(fmt.Sprintf("\tins.SetModel(models.New%sModel())\n", t.camelCase(tab)))
	b.WriteString(fmt.Sprintf("\tins.SetFormatter(New%sFormatter())\n", t.camelCase(tab)))
	b.WriteString(fmt.Sprintf("\tins.SetHooks(New%sHooks())\n", t.camelCase(tab)))
	b.WriteString(fmt.Sprintf("\treturn ins\n}"))

	return b
}

func (t *engine) createColumnFormateStr(tab string, b *strings.Builder) *strings.Builder {
	b.WriteString(fmt.Sprintf("func columnFormate(key string, value interface{}) interface{} {\n"))
	b.WriteString(fmt.Sprintf("\tkey = strings.ToLower(key)\n"))
	b.WriteString(fmt.Sprintf("\tcallback, ok := columnFuncMap[key]\n\n"))
	b.WriteString(fmt.Sprintf("\tif ok {\n"))
	b.WriteString(fmt.Sprintf("\t\treturn callback(value)\n"))
	b.WriteString(fmt.Sprintf("\t} else {\n"))
	b.WriteString(fmt.Sprintf("\t\treturn value\n"))
	b.WriteString(fmt.Sprintf("\t}\n"))
	b.WriteString(fmt.Sprintf("}\n\n\n"))

	return b
}

func (t *engine) createItemConstFile(tab string, record []columnEntry) {
	str := strings.Builder{}
	file_name := fmt.Sprintf("%s%s.go", t.ucFirst(tab), "Const")
	os.MkdirAll(fmt.Sprintf("%s/%s", t.libPath, t.ucFirst(tab)+"Dao"), 0777)
	file_path := filepath.Join(t.libPath, t.ucFirst(tab)+"Dao", file_name)

	str.WriteString(fmt.Sprintf("package %s\n\n", t.ucFirst(tab)+"Dao"))

	fmt.Println(file_path)
	err := os.WriteFile(file_path, []byte(str.String()), 0666)

	if err != nil {
		log.Println("gen code error: ", err)
	}
}

func (t *engine) createItemHooksFile(tab string, record []columnEntry) {
	str := strings.Builder{}
	file_name := fmt.Sprintf("%s%s.go", t.ucFirst(tab), "Hooks")
	os.MkdirAll(filepath.Join(t.libPath, t.ucFirst(tab)+"Dao"), 0777)
	file_path := filepath.Join(t.libPath, t.ucFirst(tab)+"Dao", file_name)

	str.WriteString(fmt.Sprintf("package %s\n\n", t.ucFirst(tab)+"Dao"))
	str.WriteString(fmt.Sprintf("import (\n"))
	str.WriteString(fmt.Sprintf("\t\"github.com/gollerxiong/gbox/components\"\n"))
	str.WriteString(fmt.Sprintf("\t\"%s\"\n", t.packageName+strings.Trim(t.pkgPath, ".")))
	//str.WriteString(fmt.Sprintf("\t\"time\"\n"))
	//str.WriteString(fmt.Sprintf("\t\"gorm.io/gorm\"\n"))
	str.WriteString(fmt.Sprintf(")\n\n"))

	str.WriteString(fmt.Sprintf("type %sHooks struct {\n", t.camelCase(tab)))
	str.WriteString(fmt.Sprintf("\tcomponents.BaseHooks\n"))
	str.WriteString(fmt.Sprintf("}\n\n"))

	str.WriteString(fmt.Sprintf("func (u *%sHooks) GetModel() *models.%sModel {\n", t.camelCase(tab), t.camelCase(tab)))
	str.WriteString(fmt.Sprintf("\treturn u.BaseHooks.GetModel().(*models.%sModel)\n", t.camelCase(tab)))
	str.WriteString(fmt.Sprintf("}\n\n"))

	str.WriteString(fmt.Sprintf("func New%sHooks() *%sHooks {\n", t.camelCase(tab), t.camelCase(tab)))
	str.WriteString(fmt.Sprintf("\tins := &%sHooks{}\n", t.camelCase(tab)))
	str.WriteString(fmt.Sprintf("\treturn ins\n"))
	str.WriteString(fmt.Sprintf("}\n\n"))

	fmt.Println(file_path)
	err := os.WriteFile(file_path, []byte(str.String()), 0666)

	if err != nil {
		log.Println("gen code error: ", err)
	}
}

func (t *engine) createItemFormaterFile(tab string, record []columnEntry) {
	str := strings.Builder{}
	file_name := fmt.Sprintf("%s%s.go", t.ucFirst(tab), "Formatter")
	os.MkdirAll(fmt.Sprintf("%s/%s", t.libPath, t.ucFirst(tab)+"Dao"), 0777)
	file_path := filepath.Join(t.libPath, t.ucFirst(tab)+"Dao", file_name)

	str.WriteString(fmt.Sprintf("package %s\n\n", t.ucFirst(tab)+"Dao"))
	str.WriteString(fmt.Sprintf("import (\n"))
	str.WriteString(fmt.Sprintf("\t\"github.com/gollerxiong/gbox/components\"\n)\n\n"))
	str.WriteString(fmt.Sprintf("type %sFormatter struct {\n", t.camelCase(tab)))
	str.WriteString(fmt.Sprintf("\tcomponents.BaseFormatter\n"))
	str.WriteString(fmt.Sprintf("}\n\n"))

	str.WriteString(fmt.Sprintf("func New%sFormatter() *%sFormatter {\n", t.camelCase(tab), t.camelCase(tab)))
	str.WriteString(fmt.Sprintf("\tres := &%sFormatter{\n", t.camelCase(tab)))
	str.WriteString(fmt.Sprintf("\t\tBaseFormatter: components.BaseFormatter{\n"))
	str.WriteString(fmt.Sprintf("\t\t\tColumnFuncMap: make(map[string]func(interface{}) interface{}),\n"))
	str.WriteString(fmt.Sprintf("\t\t},\n"))
	str.WriteString(fmt.Sprintf("\t}\n\n"))
	str.WriteString(fmt.Sprintf("\t// 注册列格式化函数\n"))
	str.WriteString(fmt.Sprintf("\t// res.ColumnFuncMap[\"sex\"] = res.sexCallback\n\n"))
	str.WriteString(fmt.Sprintf("\treturn res\n"))
	str.WriteString(fmt.Sprintf("}\n\n\n"))

	fmt.Println(file_path)

	err := os.WriteFile(file_path, []byte(str.String()), 0666)
	if err != nil {
		log.Println("gen code error: ", err)
	}
}

func (t *engine) createItemListFile(tab string, record []columnEntry) {
	str := strings.Builder{}
	file_name := fmt.Sprintf("%s%s.go", t.ucFirst(tab), "List")
	os.MkdirAll(fmt.Sprintf("%s/%s", t.libPath, t.ucFirst(tab)+"Dao"), 0777)
	file_path := filepath.Join(t.libPath, t.ucFirst(tab)+"Dao", file_name)

	str.WriteString(fmt.Sprintf("package %s\n\nimport (\n", t.ucFirst(tab)+"Dao"))
	str.WriteString(fmt.Sprintf("\t\"github.com/gollerxiong/gbox/components\"\n"))
	str.WriteString(fmt.Sprintf("\t\"%s\"\n", t.packageName+strings.Trim(t.pkgPath, ".")))
	str.WriteString(fmt.Sprintf("\t\"strings\"\n"))
	str.WriteString(fmt.Sprintf(")\n\n"))
	str.WriteString(fmt.Sprintf("type %sList struct {\n", t.camelCase(tab)))
	str.WriteString(fmt.Sprintf("\tcomponents.BaseList\n"))
	str.WriteString(fmt.Sprintf("\tField     string\n"))
	str.WriteString(fmt.Sprintf("\tModel     *models.%sModel\n", t.camelCase(tab)))
	str.WriteString(fmt.Sprintf("\tFormatter *%sListFormatter\n", t.camelCase(tab)))
	str.WriteString(fmt.Sprintf("}\n\n"))

	str.WriteString(fmt.Sprintf("func (u *%sList) BuildParams() {\n", t.camelCase(tab)))
	str.WriteString(fmt.Sprintf("\t if len(u.Ids) > 0 {\n"))
	str.WriteString(fmt.Sprintf("\t\t u.Connect = u.Connect.Where(\"id IN (?)\", u.Ids)\n"))
	str.WriteString(fmt.Sprintf("\t}\n\n"))
	str.WriteString(fmt.Sprintf("\toffset := (u.Page - 1) * u.PageSize\n"))
	str.WriteString(fmt.Sprintf("\tu.Connect = u.Connect.Offset(int(offset))\n\n"))
	str.WriteString(fmt.Sprintf("\tif len(u.Order) > 0 {\n"))
	str.WriteString(fmt.Sprintf("\t\torderStr := strings.Join(u.Order, \",\")\n"))
	str.WriteString(fmt.Sprintf("\t\torderStr = strings.Trim(orderStr, \",\")\n"))
	str.WriteString(fmt.Sprintf("\t\tu.Connect = u.Connect.Order(orderStr)\n"))
	str.WriteString(fmt.Sprintf("\t}\n\n"))
	str.WriteString(fmt.Sprintf("\tu.Connect = u.Connect.Limit(int(u.PageSize))\n"))
	str.WriteString(fmt.Sprintf("}\n\n"))

	str.WriteString(fmt.Sprintf("func (u *%sList) Find() (formateList []map[string]interface{}, total int64) {\n", t.camelCase(tab)))
	str.WriteString(fmt.Sprintf("\tu.BuildParams()\n"))
	str.WriteString(fmt.Sprintf("\tmodelList := []models.%sModel{}\n", t.camelCase(tab)))
	str.WriteString(fmt.Sprintf("\tu.Connect.Find(&modelList)\n"))
	str.WriteString(fmt.Sprintf("\tformateList = u.Formatter.SetList(modelList).SetFields(u.Field).Formate()\n\n"))
	str.WriteString(fmt.Sprintf("\tif u.PageNate {\n"))
	str.WriteString(fmt.Sprintf("\t\tu.Connect.Count(&u.Total)\n"))
	str.WriteString(fmt.Sprintf("\t}\n\n"))
	str.WriteString(fmt.Sprintf("\ttotal = u.Total\n\n"))
	str.WriteString(fmt.Sprintf("\t return\n"))
	str.WriteString(fmt.Sprintf("}\n\n"))

	str.WriteString(fmt.Sprintf("func (u *%sList) FindRaw() (modelList []models.%sModel, total int64) {\n", t.camelCase(tab), t.camelCase(tab)))
	str.WriteString(fmt.Sprintf("\tu.BuildParams()\n"))
	str.WriteString(fmt.Sprintf("\tmodelList = []models.%sModel{}\n", t.camelCase(tab)))
	str.WriteString(fmt.Sprintf("\tu.Connect.Find(&modelList)\n\n"))
	str.WriteString(fmt.Sprintf("\tif u.PageNate {\n"))
	str.WriteString(fmt.Sprintf("\t\tu.Connect.Count(&u.Total)\n"))
	str.WriteString(fmt.Sprintf("\t}\n\n"))
	str.WriteString(fmt.Sprintf("\ttotal = u.Total\n\n"))
	str.WriteString(fmt.Sprintf("\treturn\n\n"))
	str.WriteString(fmt.Sprintf("}\n\n"))

	str.WriteString(fmt.Sprintf("func New%sList() *%sList {\n", t.camelCase(tab), t.camelCase(tab)))
	str.WriteString(fmt.Sprintf("\tins := &%sList{\n", t.camelCase(tab)))
	str.WriteString(fmt.Sprintf("\t\tField:     \"*\",\n"))
	str.WriteString(fmt.Sprintf("\t\tModel:     models.New%sModel(),\n", t.camelCase(tab)))
	str.WriteString(fmt.Sprintf("\t\tFormatter: New%sListFormatter(),\n", t.camelCase(tab)))
	str.WriteString(fmt.Sprintf("\t}\n\n"))
	str.WriteString(fmt.Sprintf("\tins.Connect = ins.Model.GetConnect()\n"))
	str.WriteString(fmt.Sprintf("\tins.SetPage(1)\n"))
	str.WriteString(fmt.Sprintf("\tins.SetPageSize(20)\n"))
	str.WriteString(fmt.Sprintf("\tins.SetPageNate(true)\n"))
	str.WriteString(fmt.Sprintf("\tins.SetOrder(\"id desc\")\n\n"))
	str.WriteString(fmt.Sprintf("\treturn ins\n"))
	str.WriteString(fmt.Sprintf("}"))
	fmt.Println(file_path)

	err := os.WriteFile(file_path, []byte(str.String()), 0666)

	if err != nil {
		log.Println("gen code error: ", err)
	}
}

func (t *engine) createItemListgetFormatter(tab string, b *strings.Builder) *strings.Builder {
	b.WriteString(fmt.Sprintf("func (u *ObjectList) getFormatter() *objectListFormatter {\n"))
	b.WriteString(fmt.Sprintf("\treturn u.objectListFOrmatter\n"))
	b.WriteString(fmt.Sprintf("}\n\n"))
	t.createItemListaddField(tab, b)
	return b
}

func (t *engine) createItemListaddField(tab string, b *strings.Builder) *strings.Builder {
	b.WriteString(fmt.Sprintf("func (u *ObjectList) addField(f string) {\n"))
	b.WriteString(fmt.Sprintf("\tif !strings.Contains(u.field, \"*\") {\n"))
	b.WriteString(fmt.Sprintf("\t\tu.field += \",\" + f\n"))
	b.WriteString(fmt.Sprintf("\t}\n"))
	b.WriteString(fmt.Sprintf("}\n\n"))
	t.createItemListbuildParams(tab, b)
	return b
}

func (t *engine) createItemListbuildParams(tab string, b *strings.Builder) *strings.Builder {
	b.WriteString(fmt.Sprintf("func (u *ObjectList) buildParams() {\n"))
	b.WriteString(fmt.Sprintf("\tif len(u.ids) > 0 {\n"))
	b.WriteString(fmt.Sprintf("\t\tu.model = u.model.Where(\"id in?\", u.ids)\n"))
	b.WriteString(fmt.Sprintf("\t}\n\n"))
	b.WriteString(fmt.Sprintf("\tu.model = u.model.Offset((u.page - 1) * u.pageSize).Limit(u.pageSize)\n\n"))
	b.WriteString(fmt.Sprintf("\tif len(u.order) > 0 {\n"))
	b.WriteString(fmt.Sprintf("\t\tfor _, val := range u.order {\n"))
	b.WriteString(fmt.Sprintf("\t\t\tu.model = u.model.Order(val)\n"))
	b.WriteString(fmt.Sprintf("\t\t}\n"))
	b.WriteString(fmt.Sprintf("\t}\n"))
	b.WriteString(fmt.Sprintf("}\n\n"))
	t.createItemListFind(tab, b)
	return b
}

func (t *engine) createItemListFind(tab string, b *strings.Builder) *strings.Builder {
	b.WriteString(fmt.Sprintf("func (u *ObjectList) Find() (userInfoList []map[string]interface{}, total int64) {\n"))
	b.WriteString(fmt.Sprintf("\tu.buildParams()\n\n"))
	b.WriteString(fmt.Sprintf("\tlist := []models.%s{}\n", t.camelCase(tab)))
	b.WriteString(fmt.Sprintf("\tresult := u.model.Find(&list)\n"))
	b.WriteString(fmt.Sprintf("\tformatter := u.getFormatter()\n"))
	b.WriteString(fmt.Sprintf("\tformatter.setList(list)\n"))
	b.WriteString(fmt.Sprintf("\tformatter.setFields(u.field)\n"))
	b.WriteString(fmt.Sprintf("\tuserInfoList = formatter.formate()\n"))
	b.WriteString(fmt.Sprintf("\tu.total = result.RowsAffected\n"))
	b.WriteString(fmt.Sprintf("\ttatal = u.total\n"))
	b.WriteString(fmt.Sprintf("\treturn\n"))
	b.WriteString(fmt.Sprintf("}\n\n"))
	t.createItemListfilter(tab, b)
	return b
}

func (t *engine) createItemListfilter(tab string, b *strings.Builder) *strings.Builder {
	b.WriteString(fmt.Sprintf("func (u *ObjectList) filter(list []models.%s) (result []map[string]interface{}) {\n", t.camelCase(tab)))
	b.WriteString(fmt.Sprintf("\tfor _, v := range list {\n"))
	b.WriteString(fmt.Sprintf("\t\ttmp, _ := utils.StructToMap(v)\n"))
	b.WriteString(fmt.Sprintf("\t\ttmp = utils.PickFieldsFromMap(tmp, u.field)\n"))
	b.WriteString(fmt.Sprintf("\t\tresult = append(result, tmp)\n"))
	b.WriteString(fmt.Sprintf("\t}\n"))
	b.WriteString(fmt.Sprintf("\treturn\n"))
	b.WriteString(fmt.Sprintf("}\n\n"))
	t.createItemListObjectList(tab, b)
	return b
}

func (t *engine) createItemListObjectList(tab string, b *strings.Builder) *strings.Builder {
	b.WriteString(fmt.Sprintf("func NewObjectList() *ObjectList {\n"))
	b.WriteString(fmt.Sprintf("\treturn &ObjectList{\n"))
	b.WriteString(fmt.Sprintf("\t\tfield:\t\"*\",\n"))
	b.WriteString(fmt.Sprintf("\t\tpage:\t1,\n"))
	b.WriteString(fmt.Sprintf("\t\tpageNate:\ttrue,\n"))
	b.WriteString(fmt.Sprintf("\t\tpageSize:\t20,\n"))
	b.WriteString(fmt.Sprintf("\t\tmodel:\tGetConn().Table(table),\n"))
	b.WriteString(fmt.Sprintf("\t\torder:\t[]string{\"id desc\"},\n"))
	b.WriteString(fmt.Sprintf("\t}\n"))
	b.WriteString(fmt.Sprintf("}\n\n"))
	return b
}

func (t *engine) createItemListFormaterFile(tab string, record []columnEntry) {
	str := strings.Builder{}
	file_name := fmt.Sprintf("%s%s.go", t.ucFirst(tab), "ListFormatter")
	os.MkdirAll(fmt.Sprintf("%s/%s", t.libPath, t.ucFirst(tab)+"Dao"), 0777)
	file_path := filepath.Join(t.libPath, t.ucFirst(tab)+"Dao", file_name)

	str.WriteString(fmt.Sprintf("package %s\n\nimport (\n", t.ucFirst(tab)+"Dao"))
	str.WriteString(fmt.Sprintf("\t\"github.com/gollerxiong/gbox/components\"\n"))
	str.WriteString(fmt.Sprintf("\t\"%s\"\n", t.packageName+strings.Trim(t.pkgPath, ".")))
	str.WriteString(fmt.Sprintf("\t\"strings\"\n"))
	str.WriteString(fmt.Sprintf("\t\"sync\"\n"))
	str.WriteString(fmt.Sprintf(")\n\n"))

	str.WriteString(fmt.Sprintf("type %sListFormatter struct {\n", t.camelCase(tab)))
	str.WriteString(fmt.Sprintf("\tcomponents.BaseListFormatter\n"))
	str.WriteString(fmt.Sprintf("\tList      []models.%sModel\n", t.camelCase(tab)))
	str.WriteString(fmt.Sprintf("\tFormatter *%sFormatter\n", t.camelCase(tab)))
	str.WriteString(fmt.Sprintf("\tmutex     sync.Mutex\n"))
	str.WriteString(fmt.Sprintf("}\n\n"))

	str.WriteString(fmt.Sprintf("func (u *%sListFormatter) SetList(list []models.%sModel) *%sListFormatter {\n", t.camelCase(tab), t.camelCase(tab), t.camelCase(tab)))
	str.WriteString(fmt.Sprintf("\tu.List = list\n"))
	str.WriteString(fmt.Sprintf("\treturn u\n"))
	str.WriteString(fmt.Sprintf("}\n\n"))

	str.WriteString(fmt.Sprintf("func (u *%sListFormatter) SetFields(field string) *%sListFormatter {\n", t.camelCase(tab), t.camelCase(tab)))
	str.WriteString(fmt.Sprintf("\tu.Fields = field\n"))
	str.WriteString(fmt.Sprintf("\treturn u\n"))
	str.WriteString(fmt.Sprintf("}\n\n"))

	str.WriteString(fmt.Sprintf("func (u *%sListFormatter) Formate() []map[string]interface{} {\n", t.camelCase(tab)))
	str.WriteString(fmt.Sprintf("\tvar length = len(u.List)\n\n"))
	str.WriteString(fmt.Sprintf("\tu.WG.Add(length)\n"))
	str.WriteString(fmt.Sprintf("\tfor i := 0; i < length; i++ {\n"))
	str.WriteString(fmt.Sprintf("\t\tgo u.Do(u.List[i])\n"))
	str.WriteString(fmt.Sprintf("\t}\n\n"))
	str.WriteString(fmt.Sprintf("\tu.WG.Wait()\n"))
	str.WriteString(fmt.Sprintf("\treturn u.Result\n"))
	str.WriteString(fmt.Sprintf("}\n\n"))

	str.WriteString(fmt.Sprintf("func (b *%sListFormatter) Do(item models.%sModel) {\n", t.camelCase(tab), t.camelCase(tab)))
	str.WriteString(fmt.Sprintf("\tdefer b.WG.Done()\n"))
	str.WriteString(fmt.Sprintf("\tresult := make(map[string]interface{})\n"))
	str.WriteString(fmt.Sprintf("\tarr := components.StructToMap(item)\n"))
	str.WriteString(fmt.Sprintf("\tif !strings.Contains(b.Fields, \"*\") {\n"))
	str.WriteString(fmt.Sprintf("\t\tfieldArr := strings.Split(b.Fields, \",\")\n"))
	str.WriteString(fmt.Sprintf("\t\tfor _, key := range fieldArr {\n"))
	str.WriteString(fmt.Sprintf("\t\t\tval, ok := arr[key]\n"))
	str.WriteString(fmt.Sprintf("\t\t\tif ok {\n"))
	str.WriteString(fmt.Sprintf("\t\t\t\tresult[key] = b.Formatter.ColumnFormate(key, val)\n"))
	str.WriteString(fmt.Sprintf("\t\t\t}\n"))
	str.WriteString(fmt.Sprintf("\t\t}\n"))
	str.WriteString(fmt.Sprintf("\t\tb.mutex.Lock() // 加锁\n"))
	str.WriteString(fmt.Sprintf("\t\tb.Result = append(b.Result, result)\n"))
	str.WriteString(fmt.Sprintf("\t\tb.mutex.Unlock() // 解锁\n"))
	str.WriteString(fmt.Sprintf("\t} else {\n"))
	str.WriteString(fmt.Sprintf("\t\tb.mutex.Lock() // 加锁\n"))
	str.WriteString(fmt.Sprintf("\t\tb.Result = append(b.Result, arr)\n"))
	str.WriteString(fmt.Sprintf("\t\tb.mutex.Unlock() // 解锁\n"))
	str.WriteString(fmt.Sprintf("\t}\n"))
	str.WriteString(fmt.Sprintf("}\n\n"))

	str.WriteString(fmt.Sprintf("func New%sListFormatter() *%sListFormatter {\n", t.camelCase(tab), t.camelCase(tab)))
	str.WriteString(fmt.Sprintf("\treturn &%sListFormatter{\n", t.camelCase(tab)))
	str.WriteString(fmt.Sprintf("\t\tFormatter: New%sFormatter(),\n", t.camelCase(tab)))
	str.WriteString(fmt.Sprintf("\t}\n"))
	str.WriteString(fmt.Sprintf("}\n"))

	fmt.Println(file_path)
	err := os.WriteFile(file_path, []byte(str.String()), 0666)

	if err != nil {
		log.Println("gen code error: ", err)
	}
}

func (t *engine) createItemListnewObjectListFormatterStr(tab string, b *strings.Builder) *strings.Builder {
	b.WriteString(fmt.Sprintf("type %sListFormatter struct {\n", t.camelCase(tab)))
	b.WriteString(fmt.Sprintf("\tcomponents.BaseListFormatter\n"))
	b.WriteString(fmt.Sprintf("\tList      []models.%sModel\n", t.camelCase(tab)))
	b.WriteString(fmt.Sprintf("\tFormatter *%sFormatter\n", t.camelCase(tab)))
	b.WriteString(fmt.Sprintf("\tmutex     sync.Mutex\n"))
	b.WriteString(fmt.Sprintf("}\n\n"))
	t.createItemListobjectListFormatter(tab, b)
	return b
}

func (t *engine) createItemListobjectListFormatter(tab string, b *strings.Builder) *strings.Builder {
	b.WriteString(fmt.Sprintf("type objectListFormatter struct {\n"))
	b.WriteString(fmt.Sprintf("\tlist\t[]map[string]interface{}\n"))
	b.WriteString(fmt.Sprintf("\tfields\tstring\n"))
	b.WriteString(fmt.Sprintf("\twg\tsync.WaitGroup\n"))
	b.WriteString(fmt.Sprintf("\tresult\t[]map[string]interface{}\n"))
	b.WriteString(fmt.Sprintf("}\n\n"))
	t.createItemListsetList(tab, b)
	return b
}

func (t *engine) createItemListsetList(tab string, b *strings.Builder) *strings.Builder {
	b.WriteString(fmt.Sprintf("func (u *objectListFormatter) setList(list []models.%s) {\n", t.camelCase(tab)))
	b.WriteString(fmt.Sprintf("\tvar end = len(list)\n"))
	b.WriteString(fmt.Sprintf("\tvar tmp = make([]map[string]interface{}, end)\n\n"))
	b.WriteString(fmt.Sprintf("\tu.result = tmp\n\n"))
	b.WriteString(fmt.Sprintf("\tfor i := 0; i < end; i++ {\n"))
	b.WriteString(fmt.Sprintf("\t\ttmpmap, _ := utils.StructToMap(list[i])\n"))
	b.WriteString(fmt.Sprintf("\t\ttmp[i] = tmpmap\n"))
	b.WriteString(fmt.Sprintf("\t}\n"))
	b.WriteString(fmt.Sprintf("\tu.list = tmp\n"))
	b.WriteString(fmt.Sprintf("}\n\n"))
	t.createItemListsetFields(tab, b)
	return b
}

func (t *engine) createItemListsetFields(tab string, b *strings.Builder) *strings.Builder {
	b.WriteString(fmt.Sprintf("func (u *objectListFormatter) setFields(f string) {\n"))
	b.WriteString(fmt.Sprintf("\tu.fields = f\n"))
	b.WriteString(fmt.Sprintf("}\n\n"))

	b.WriteString(fmt.Sprintf("func (u *objectListFormatter) formate() []map[string]interface{} {\n"))
	b.WriteString(fmt.Sprintf("\tvar length = len(u.list)\n\n"))
	b.WriteString(fmt.Sprintf("\tu.wg.Add(length)\n"))
	b.WriteString(fmt.Sprintf("\tfor i := 0; i < length; i++ {\n"))
	b.WriteString(fmt.Sprintf("\t\tgo u.do(i, u.list[i])\n"))
	b.WriteString(fmt.Sprintf("\t}\n\n"))
	b.WriteString(fmt.Sprintf("\tu.wg.Wait()\n"))
	b.WriteString(fmt.Sprintf("\treturn u.result\n"))
	b.WriteString(fmt.Sprintf("}\n\n"))

	b.WriteString(fmt.Sprintf("func (u *objectListFormatter) do(index int, item map[string]interface{}) {\n"))
	b.WriteString(fmt.Sprintf("\tdefer u.wg.Done()\n"))
	b.WriteString(fmt.Sprintf("\tresult := make(map[string]interface{})\n\n"))
	b.WriteString(fmt.Sprintf("\tif !strings.Contains(u.fields, \"*\") {\n"))
	b.WriteString(fmt.Sprintf("\t\tfieldArr := strings.Split(u.fields, \",\")\n"))
	b.WriteString(fmt.Sprintf("\t\tfor _, key := range fieldArr {\n"))
	b.WriteString(fmt.Sprintf("\t\t\tval, ok := item[key]\n"))
	b.WriteString(fmt.Sprintf("\t\t\tif ok {\n"))
	b.WriteString(fmt.Sprintf("\t\t\t\tresult[key] = columnFormate(key, val)\n"))
	b.WriteString(fmt.Sprintf("\t\t\t}\n"))
	b.WriteString(fmt.Sprintf("\t\t}\n"))
	b.WriteString(fmt.Sprintf("\t\tu.result[index] = result\n"))
	b.WriteString(fmt.Sprintf("\t} else {\n"))
	b.WriteString(fmt.Sprintf("\t\tu.result[index] = item\n"))
	b.WriteString(fmt.Sprintf("\t}\n"))
	b.WriteString(fmt.Sprintf("}\n\n\n"))
	return b
}

func (t *engine) createItemHelperFile(tab string, record []columnEntry) {
	str := strings.Builder{}
	file_name := fmt.Sprintf("%s%s.go", t.ucFirst(tab), "Helper")
	os.MkdirAll(fmt.Sprintf("%s/%s", t.libPath, t.ucFirst(tab)+"Dao"), 0777)
	file_path := filepath.Join(t.libPath, t.ucFirst(tab)+"Dao", file_name)

	str.WriteString(fmt.Sprintf("package %s\n\n", t.ucFirst(tab)+"Dao"))
	str.WriteString(fmt.Sprintf("type %sHelper struct {}\n\n", t.camelCase(tab)))
	str.WriteString(fmt.Sprintf("func New%sHelper() *%sHelper {\n", t.camelCase(tab), t.camelCase(tab)))
	str.WriteString(fmt.Sprintf("\treturn &%sHelper{}\n", t.camelCase(tab)))
	str.WriteString(fmt.Sprintf("}\n"))

	fmt.Println(file_path)
	err := os.WriteFile(file_path, []byte(str.String()), 0666)

	if err != nil {
		log.Println("gen code error: ", err)
	}
}

func (t *engine) createItemBatchOperatorFile(tab string, record []columnEntry) {
	str := strings.Builder{}
	file_name := fmt.Sprintf("%s%s.go", t.ucFirst(tab), "BatchOperator")
	os.MkdirAll(fmt.Sprintf("%s/%s", t.libPath, t.ucFirst(tab)+"Dao"), 0777)
	file_path := filepath.Join(t.libPath, t.ucFirst(tab)+"Dao", file_name)

	str.WriteString(fmt.Sprintf("package %s\n\nimport (\n", t.ucFirst(tab)+"Dao"))
	str.WriteString(fmt.Sprintf("\t\"%s/components\"\n", "github.com/gollerxiong/gbox"))
	str.WriteString(fmt.Sprintf("\t\"%s\"\n", t.packageName+strings.Trim(t.pkgPath, ".")))
	str.WriteString(fmt.Sprintf(")\n\n"))

	str.WriteString(fmt.Sprintf("type %sBatchOperator struct {\n", t.camelCase(tab)))
	str.WriteString(fmt.Sprintf("\tcomponents.BaseBatchOperator\n"))
	str.WriteString(fmt.Sprintf("}\n\n"))

	str.WriteString(fmt.Sprintf("func (b *%sBatchOperator) buildParams() {\n", t.camelCase(tab)))
	str.WriteString(fmt.Sprintf("\tif len(b.Ids) > 0 {\n"))
	str.WriteString(fmt.Sprintf("\t\tb.Connect = b.Connect.Where(\"id in ?\", b.Ids)\n"))
	str.WriteString(fmt.Sprintf("\t}\n"))
	str.WriteString(fmt.Sprintf("}\n\n"))

	str.WriteString(fmt.Sprintf("func (b *%sBatchOperator) Update() bool {\n", t.camelCase(tab)))
	str.WriteString(fmt.Sprintf("\tb.buildParams()\n"))
	str.WriteString(fmt.Sprintf("\tb.Connect.Model(&models.%sModel{}).Update(b.Field, b.FieldValue)\n", t.camelCase(tab)))
	str.WriteString(fmt.Sprintf("\treturn true\n"))
	str.WriteString(fmt.Sprintf("}\n\n"))

	str.WriteString(fmt.Sprintf("func (b *%sBatchOperator) Delete() bool {\n", t.camelCase(tab)))
	str.WriteString(fmt.Sprintf("\tb.buildParams()\n"))
	str.WriteString(fmt.Sprintf("\tb.Connect.Delete(&models.%sModel{})\n", t.camelCase(tab)))
	str.WriteString(fmt.Sprintf("\treturn true\n"))
	str.WriteString(fmt.Sprintf("}\n\n"))

	str.WriteString(fmt.Sprintf("func New%sBatchOperator() *%sBatchOperator {\n", t.camelCase(tab), t.camelCase(tab)))
	str.WriteString(fmt.Sprintf("\tins := &%sBatchOperator{}\n", t.camelCase(tab)))
	str.WriteString(fmt.Sprintf("\tins.SetModel(models.New%sModel())\n", t.camelCase(tab)))
	str.WriteString(fmt.Sprintf("\tins.Connect = ins.Model.GetConnect()\n"))
	str.WriteString(fmt.Sprintf("\treturn ins\n"))
	str.WriteString(fmt.Sprintf("}\n\n"))

	fmt.Println(file_path)

	err := os.WriteFile(file_path, []byte(str.String()), 0666)

	if err != nil {
		log.Println("gen code error: ", err)
	}
}

func (t *engine) createLibrary(tab string, record []columnEntry) {
	t.createItemFile(tab, record)
	t.createItemFormaterFile(tab, record)
	t.createItemListFile(tab, record)
	t.createItemListFormaterFile(tab, record)
	t.createItemHelperFile(tab, record)
	t.createItemBatchOperatorFile(tab, record)
	t.createItemConstFile(tab, record)
	t.createItemHooksFile(tab, record)
}

func checkPathExist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}

	if os.IsNotExist(err) {
		return false
	}

	return false
}

func strInSet(s string, sets []string) bool {
	for _, member := range sets {
		if s == member {
			return true
		}
	}

	return false
}
