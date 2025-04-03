package components

import (
	"github.com/gollerxiong/gbox/components/interfaces"
	"gorm.io/gorm"
)

type BaseHooks struct {
	model interfaces.ModelInterface
}

func (b *BaseHooks) BeforeSave(*gorm.DB)   {}
func (b *BaseHooks) BeforeCreate(*gorm.DB) {}
func (b *BaseHooks) AfterCreate(*gorm.DB)  {}
func (b *BaseHooks) AfterSave(*gorm.DB)    {}

func (b *BaseHooks) BeforeUpdate(*gorm.DB) {}
func (b *BaseHooks) AfterUpdate(*gorm.DB)  {}
func (b *BaseHooks) BeforeDelete(*gorm.DB) {}
func (b *BaseHooks) AfterDelete(*gorm.DB)  {}
func (b *BaseHooks) AfterFind(*gorm.DB)    {}
func (b *BaseHooks) SetModel(model interfaces.ModelInterface) interface{} {
	b.model = model
	return b
}
func (b *BaseHooks) GetModel() interface{} {
	return b.model
}
