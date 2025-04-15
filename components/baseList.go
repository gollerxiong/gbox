package components

import (
	"github.com/gollerxiong/gbox/components/interfaces"
	"gorm.io/gorm"
)

type BaseList struct {
	Field     string
	Model     *interfaces.ModelInterface
	Connect   *gorm.DB
	Page      int64
	PageNate  bool
	PageSize  int64
	Order     []string
	Ids       []int64
	List      []interface{}
	Total     int64
	Formatter *interfaces.ListFormatterInterface
}

func (b *BaseList) GetFormatter() interface{} {
	return b.Formatter
}

func (b *BaseList) SetField(field string) *BaseList {
	b.Field = field
	return b
}

func (b *BaseList) SetModel(model *interfaces.ModelInterface) *BaseList {
	b.Model = model
	return b
}
func (b *BaseList) SetPage(page int64) *BaseList {
	b.Page = page
	return b
}

func (b *BaseList) SetPageNate(pageNate bool) *BaseList {
	b.PageNate = pageNate
	return b
}

func (b *BaseList) SetPageSize(pageSize int64) *BaseList {
	b.PageSize = pageSize
	return b
}

func (b *BaseList) SetOrder(order string) *BaseList {
	b.Order = append(b.Order, order)
	return b
}

func (b *BaseList) SetIds(ids []int64) *BaseList {
	b.Ids = ids
	return b
}

func (b *BaseList) BuildParams() {
	if len(b.Ids) > 0 {
		b.Connect = b.Connect.Where("id IN (?)", b.Ids)
	}

	offset := (b.Page - 1) * b.PageSize
	b.Connect = b.Connect.Offset(int(offset)).Limit(int(b.PageSize))
}
