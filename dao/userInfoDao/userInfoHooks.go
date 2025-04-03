package userInfoDao

import (
	"github.com/gollerxiong/gbox/components"
	"github.com/gollerxiong/gbox/models"
	"time"

	"gorm.io/gorm"
)

type UserInfoHooks struct {
	components.BaseHooks
}

func (u *UserInfoHooks) GetModel() *models.UserInfoModel {
	return u.BaseHooks.GetModel().(*models.UserInfoModel)
}

func (u *UserInfoHooks) BeforeSave(db *gorm.DB) {
	currentTime := time.Now()
	// 将当前时间格式化为字符串
	timeStr := currentTime.Format("2006-01-02 15:04:05")
	u.GetModel().Nickname = timeStr + "_save"
}

func (u *UserInfoHooks) BeforeUpdate(db *gorm.DB) {
	currentTime := time.Now()
	// 将当前时间格式化为字符串
	timeStr := currentTime.Format("2006-01-02 15:04:05")
	u.GetModel().Nickname = timeStr + "_update"
}

//func (u *UserInfoHooks) SetModel(model *models.UserInfoModel) interface{} {
//	u.model = model
//	return u
//}

func NewUserInfoHooks() *UserInfoHooks {
	ins := &UserInfoHooks{}
	//ins.SetModel(models.NewUserInfoModel())
	return ins
}
