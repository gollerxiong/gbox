package userInfoDao

import (
	"github.com/gollerxiong/gbox/components"
	"github.com/gollerxiong/gbox/models"
)

type UserInfoBatchOperator struct {
	components.BaseBatchOperator
}

func (b *UserInfoBatchOperator) buildParams() {
	if len(b.Ids) > 0 {
		b.Connect = b.Connect.Where("id in ?", b.Ids)
	}
}

func (b *UserInfoBatchOperator) Update() bool {
	b.buildParams()
	b.Connect.Model(&models.UserInfoModel{}).Update(b.Field, b.FieldValue)
	return true
}

func (b *UserInfoBatchOperator) Delete() bool {
	b.buildParams()
	b.Connect.Delete(&models.UserInfoModel{})
	return true
}

func NewUserInfoBatchOperator() *UserInfoBatchOperator {
	ins := &UserInfoBatchOperator{}
	ins.SetModel(models.NewUserInfoModel())
	ins.Connect = ins.Model.GetConnect()
	return ins
}
