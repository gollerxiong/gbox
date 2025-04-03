package userInfoExtendDao

import (
	"gormlibrary/components"
)

type UserInfoExtendFormatter struct {
	components.BaseFormatter
}

func NewUserInfoExtendFormatter() *UserInfoExtendFormatter {
	res := &UserInfoExtendFormatter{
		BaseFormatter: components.BaseFormatter{
			ColumnFuncMap: make(map[string]func(interface{}) interface{}),
		},
	}

	return res
}
