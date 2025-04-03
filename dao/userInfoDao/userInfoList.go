package userInfoDao

import (
	"github.com/gollerxiong/gbox/components"
	"github.com/gollerxiong/gbox/models"
	"strings"
)

type UserInfoList struct {
	components.BaseList
	Field     string
	Model     *models.UserInfoModel
	Formatter *UserInfoListFormatter
}

func (u *UserInfoList) BuildParams() {
	if len(u.Ids) > 0 {
		u.Connect = u.Connect.Where("id IN (?)", u.Ids)
	}

	offset := (u.Page - 1) * u.PageSize
	u.Connect = u.Connect.Offset(int(offset))

	if len(u.Order) > 0 {
		orderStr := strings.Join(u.Order, ",")
		orderStr = strings.Trim(orderStr, ",")
		u.Connect = u.Connect.Order(orderStr)
	}

	u.Connect = u.Connect.Limit(int(u.PageSize))
}

func (u *UserInfoList) Find() (formateList []map[string]interface{}, total int64) {
	u.BuildParams()
	modelList := []models.UserInfoModel{}
	u.Connect.Find(&modelList)
	formateList = u.Formatter.SetList(modelList).SetFields(u.Field).Formate()

	if u.PageNate {
		u.Connect.Count(&u.Total)
	}

	total = u.Total

	return
}

func (u *UserInfoList) FindRaw() (modelList []models.UserInfoModel, total int64) {
	u.BuildParams()
	u.Connect.Find(&modelList)

	if u.PageNate {
		u.Connect.Count(&u.Total)
	}

	total = u.Total

	return
}

func NewUserInfoList() *UserInfoList {
	ins := &UserInfoList{
		Field:     "*",
		Model:     models.NewUserInfoModel(),
		Formatter: NewUserInfoListFormatter(),
	}

	ins.Connect = ins.Model.GetConnect()
	ins.SetPage(1)
	ins.SetPageSize(20)
	ins.SetPageNate(true)
	ins.SetOrder("id desc")

	return ins
}
