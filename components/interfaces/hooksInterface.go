package interfaces

import "gorm.io/gorm"

type HooksInterface interface {
	BeforeSave(*gorm.DB)
	BeforeCreate(*gorm.DB)
	AfterCreate(*gorm.DB)
	AfterSave(*gorm.DB)

	BeforeUpdate(*gorm.DB)
	AfterUpdate(*gorm.DB)
	BeforeDelete(*gorm.DB)
	AfterDelete(*gorm.DB)
	AfterFind(*gorm.DB)

	SetModel(ModelInterface) interface{}
}
