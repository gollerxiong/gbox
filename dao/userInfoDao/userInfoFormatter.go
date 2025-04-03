package userInfoDao

import (
	"github.com/gollerxiong/gbox/components"
)

type UserInfoFormatter struct {
	components.BaseFormatter
}

func (u *UserInfoFormatter) sexCallback(sex interface{}) interface{} {
	if sex.(int64) == 0 {
		return "男"
	} else {
		return "女"
	}
}

func NewUserInfoFormatter() *UserInfoFormatter {
	res := &UserInfoFormatter{
		BaseFormatter: components.BaseFormatter{
			ColumnFuncMap: make(map[string]func(interface{}) interface{}),
		},
	}

	// 注册列格式化函数
	res.ColumnFuncMap["sex"] = res.sexCallback

	return res
}
