package components

import (
	"github.com/gollerxiong/gbox/components/interfaces"
	"gorm.io/gorm"
)

type BaseObject struct {
	isNew         bool
	params        map[string]interface{}
	model         interfaces.ModelInterface
	oldModel      interfaces.ModelInterface
	field         string
	formatter     interfaces.FormatterInterface
	attributes    map[string]interface{}
	oldAttributes map[string]interface{}
	hooks         interfaces.HooksInterface
}

func (b *BaseObject) GetBaseObject() interface{} {
	return b
}

func (b *BaseObject) GetModel() interface{} {
	return b.model
}

func (b *BaseObject) GetFormatter() interface{} {
	return b.formatter
}

func (b *BaseObject) GetAttributes() map[string]interface{} {
	// 筛选出需要的字段
	result := make(map[string]interface{})

	if !strings.Contains(b.field, "*") {
		fieldArr := strings.Split(b.field, ",")

		for _, key := range fieldArr {
			val, ok := b.attributes[key]

			if ok {
				result[key] = val
			}
		}
	} else {
		result = b.attributes
	}

	// 对字段做格式化
	b.formatter.SetData(result)

	result = b.formatter.Formate()

	return result
}

func (b *BaseObject) GetOldAttributes() map[string]interface{} {
	// 筛选出需要的字段
	result := make(map[string]interface{})

	if !strings.Contains(b.field, "*") {
		fieldArr := strings.Split(b.field, ",")

		for _, key := range fieldArr {
			val, ok := b.oldAttributes[key]

			if ok {
				result[key] = val
			}
		}
	} else {
		result = b.attributes
	}

	// 对字段做格式化
	b.formatter.SetData(result)

	result = b.formatter.Formate()

	return result
}

func (b *BaseObject) GetHooks() interface{} {
	return b.hooks
}

func (b *BaseObject) IsNew() bool {
	return b.isNew
}

func (b *BaseObject) Save() bool {
	b.SetAttributes(StructToMap(b.model))
	hooks := b.GetHooks().(interfaces.HooksInterface)
	hooks.SetModel(b.model)

	if b.IsNew() {
		tx := b.model.GetConnect().Begin()
		defer func() {
			if err := recover(); err != nil {
				tx.Rollback()
				panic(err)
			}
		}()

		hooks.BeforeSave(b.model.GetConnect())
		result := tx.Create(b.model)
		if result.Error != nil {
			panic(result.Error.Error())
		}
		hooks.AfterSave(b.model.GetConnect())
		tx.Commit()
	} else {
		hooks.BeforeUpdate(b.model.GetConnect())
		b.SetAttributes(StructToMap(b.model))
		diffMap := DiffMapBaseFirst(b.GetOldAttributes(), b.GetAttributes())
		keys := MapKeys(diffMap)
		tx := b.model.GetConnect().Begin()
		defer func() {
			if err := recover(); err != nil {
				tx.Rollback()
				panic(err)
			}
		}()

		result := tx.Model(b.GetModel().(interfaces.ModelInterface)).Where(b.model.GetPrimaryKey()+" = ?", b.model.GetId()).Select(keys).Updates(diffMap)
		err := result.Error
		if err != nil {
			panic(err)
		}
		hooks.AfterUpdate(b.model.GetConnect())
		tx.Commit()
	}

	b.SetNew(false)
	b.SetAttributes(StructToMap(b.model))
	return true
}

func (b *BaseObject) Delete() bool {
	model := *b.GetModel().(*interfaces.ModelInterface)
	if b.IsNew() || model.GetId() < 1 {
		return true
	} else {
		res := model.GetConnect().Delete(b.model).Error

		if res != nil {
			return false
		}
	}

	return true
}

func (b *BaseObject) LoadById(id int64) interface{} {
	model := b.GetModel().(interfaces.ModelInterface)
	err := model.GetConnect().Where(b.model.GetPrimaryKey()+" = ?", id).First(model).Error

	if err == gorm.ErrRecordNotFound {
		b.SetNew(true)
	} else {
		b.SetNew(false)
		b.SetModel(model)
		b.SetAttributes(StructToMap(model))
		b.SetOldAttributes(StructToMap(model))
	}

	return b
}

func (b *BaseObject) DeleteById(id int64) bool {
	return b.LoadById(id).(*BaseObject).Delete()
}

func (b *BaseObject) SetNew(value bool) *BaseObject {
	b.isNew = value
	return b
}

func (b *BaseObject) SetField(field string) *BaseObject {
	b.field = field
	return b
}

func (b *BaseObject) SetAttributes(attributes map[string]interface{}) *BaseObject {
	b.attributes = attributes
	return b
}

func (b *BaseObject) SetOldAttributes(oldAttributes map[string]interface{}) *BaseObject {
	b.oldAttributes = oldAttributes
	return b
}

func (b *BaseObject) SetParams(params map[string]interface{}) *BaseObject {
	b.params = params
	return b
}

func (b *BaseObject) SetFormatter(formatter interfaces.FormatterInterface) *BaseObject {
	b.formatter = formatter
	return b
}

func (b *BaseObject) SetModel(model interfaces.ModelInterface) *BaseObject {
	b.model = model
	return b
}

func (b *BaseObject) SetOldModel(model interfaces.ModelInterface) *BaseObject {
	b.oldModel = model
	return b
}

func (b *BaseObject) SetHooks(hooks interfaces.HooksInterface) *BaseObject {
	b.hooks = hooks
	return b
}
