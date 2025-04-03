package userInfoExtendDao

import (
	"gorm.io/gorm"
	"gormlibrary/components"
	"gormlibrary/models"
)

type userInfoExten struct {
	components.BaseObject
}

func (u *userInfoExten) GetModel() *models.UserInfoModel {
	return u.BaseObject.GetModel().(*models.UserInfoModel)
}

func (u *userInfoExten) LoadById(id int64) *userInfoExten {
	u.BaseObject.LoadById(id)
	return u
}

func (u *userInfoExten) LoadByUserId(userId int64) *userInfoExten {
	model := u.GetModel()
	connect := model.GetConnect()
	err := connect.Where("user_id = ?", userId).First(u.GetModel()).Error

	if err == gorm.ErrRecordNotFound {
		u.SetNew(true)
	} else {
		u.SetNew(false)
		u.SetModel(model)
		u.SetAttributes(components.StructToMap(model))
		u.SetOldAttributes(components.StructToMap(model))
	}

	return u
}

/*
* 实例化对象并注入依赖
 */
func NewUserInfoExten() *userInfoExten {
	ins := &userInfoExten{}
	ins.SetField("*")
	ins.SetModel(models.NewUserInfoModel())
	ins.SetFormatter(NewUserInfoExtendFormatter())
	ins.SetHooks(NewUserInfoExtendHooks())
	return ins
}
