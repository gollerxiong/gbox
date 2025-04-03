package userInfoDao

import (
	"github.com/gollerxiong/gbox/components"
	"github.com/gollerxiong/gbox/models"
	"strings"
	"sync" // 导入sync包
)

type UserInfoListFormatter struct {
	components.BaseListFormatter
	List      []models.UserInfoModel
	Formatter *UserInfoFormatter
	mutex     sync.Mutex // 定义互斥锁
}

func (u *UserInfoListFormatter) SetList(list []models.UserInfoModel) *UserInfoListFormatter {
	u.List = list
	return u
}

func (u *UserInfoListFormatter) SetFields(field string) *UserInfoListFormatter {
	u.Fields = field
	return u
}

func (u *UserInfoListFormatter) Formate() []map[string]interface{} {
	var length = len(u.List)

	u.WG.Add(length)
	for i := 0; i < length; i++ {
		go u.Do(u.List[i])
	}

	u.WG.Wait()
	return u.Result
}

func (b *UserInfoListFormatter) Do(item models.UserInfoModel) {
	defer b.WG.Done()
	result := make(map[string]interface{})
	arr := components.StructToMap(item)
	if !strings.Contains(b.Fields, "*") {
		fieldArr := strings.Split(b.Fields, ",")
		for _, key := range fieldArr {
			val, ok := arr[key]
			if ok {
				result[key] = b.Formatter.ColumnFormate(key, val)
			}
		}
		b.mutex.Lock() // 加锁
		b.Result = append(b.Result, result)
		b.mutex.Unlock() // 解锁
	} else {
		b.mutex.Lock() // 加锁
		b.Result = append(b.Result, arr)
		b.mutex.Unlock() // 解锁
	}
}

func NewUserInfoListFormatter() *UserInfoListFormatter {
	return &UserInfoListFormatter{
		Formatter: NewUserInfoFormatter(),
	}
}
