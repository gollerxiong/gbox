package interfaces

import "gorm.io/gorm"

type ModelInterface interface {
	GetTableName() string
	GetConnectName() string
	GetConnect() *gorm.DB
	GetId() int64
}
