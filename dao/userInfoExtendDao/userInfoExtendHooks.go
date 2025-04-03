package userInfoExtendDao

import (
	"gormlibrary/components"
	"gormlibrary/models"
	"time"

	"gorm.io/gorm"
)

type UserInfoExtendHooks struct {
	components.BaseHooks
}

func (u *UserInfoExtendHooks) GetModel() *models.UserInfoModel {
	return u.BaseHooks.GetModel().(*models.UserInfoModel)
}

func (u *UserInfoExtendHooks) BeforeSave(db *gorm.DB) {
	currentTime := time.Now()
	// 将当前时间格式化为字符串
	timeStr := currentTime.Format("2006-01-02 15:04:05")
	u.GetModel().Nickname = timeStr + "_save"
}

func (u *UserInfoExtendHooks) BeforeUpdate(db *gorm.DB) {
	currentTime := time.Now()
	// 将当前时间格式化为字符串
	timeStr := currentTime.Format("2006-01-02 15:04:05")
	u.GetModel().Nickname = timeStr + "_update"
}

//func (u *UserInfoHooks) SetModel(model *models.UserInfoModel) interface{} {
//	u.model = model
//	return u
//}

func NewUserInfoExtendHooks() *UserInfoExtendHooks {
	ins := &UserInfoExtendHooks{}
	//ins.SetModel(models.NewUserInfoModel())
	return ins
}
