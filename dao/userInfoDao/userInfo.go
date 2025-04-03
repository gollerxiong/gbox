package userInfoDao

import (
	"codeup.aliyun.com/gmax/components"
	"codeup.aliyun.com/gmax/models"
	"gorm.io/gorm"
)

type UserInfo struct {
	components.BaseObject
	withExtend bool
}

func (u *UserInfo) GetModel() *models.UserInfoModel {
	return u.BaseObject.GetModel().(*models.UserInfoModel)
}

func (u *UserInfo) SetWithExtend(withExtend bool) *UserInfo {
	u.withExtend = withExtend
	return u
}

func (u *UserInfo) LoadById(id int64) *UserInfo {
	connect := u.GetModel().Connect
	// 加载用户的扩展信息
	if u.withExtend {
		connect = connect.Preload("UserInfoExtend")
	}
	u.BaseObject.LoadById(id)
	return u
}

func (u *UserInfo) LoadByUsername(username string) *UserInfo {
	model := u.GetModel()
	connect := model.GetConnect()
	err := connect.Where("user_name = ?", username).First(u.GetModel()).Error

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

/*uytuy6
* 实例化对象并注入依赖
 */
func NewUserInfo() *UserInfo {
	ins := &UserInfo{}
	ins.SetField("*")
	ins.SetModel(models.NewUserInfoModel())
	ins.SetFormatter(NewUserInfoFormatter())
	ins.SetHooks(NewUserInfoHooks())
	return ins
}
