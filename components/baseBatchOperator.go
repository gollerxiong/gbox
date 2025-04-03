package components

import (
	"github.com/gollerxiong/gbox/components/interfaces"
	"gorm.io/gorm"
)

type BaseBatchOperator struct {
	Model      interfaces.ModelInterface
	Connect    *gorm.DB
	Field      string
	FieldValue interface{}
	Ids        []int64
	Page       int64
	PageSize   int64
}

func (b *BaseBatchOperator) GetModel() interface{} {
	return b.Model
}

func (b *BaseBatchOperator) SetModel(model interfaces.ModelInterface) *BaseBatchOperator {
	b.Model = model
	return b
}

func (b *BaseBatchOperator) SetPage(page int64) *BaseBatchOperator {
	b.Page = page
	return b
}

func (b *BaseBatchOperator) SetPageSize(pageSize int64) *BaseBatchOperator {
	b.PageSize = pageSize
	return b
}

func (b *BaseBatchOperator) SetField(f string) *BaseBatchOperator {
	b.Field = f
	return b
}

func (b *BaseBatchOperator) SetIds(ids []int64) *BaseBatchOperator {
	b.Ids = ids
	return b
}

func (b *BaseBatchOperator) SetFieldValue(v interface{}) *BaseBatchOperator {
	b.FieldValue = v
	return b
}

func (b *BaseBatchOperator) BuildParams() {
	if len(b.Ids) > 0 {
		b.Connect = b.Connect.Where("id IN (?)", b.Ids)
	}

	offset := (b.Page - 1) * b.PageSize
	b.Connect = b.Connect.Offset(int(offset)).Limit(int(b.PageSize))
}
