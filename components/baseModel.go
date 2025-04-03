package components

import (
	"gorm.io/gorm"
	"time"
)

type BaseModel struct {
	ID        int64     `json:"id" gorm:"primarykey"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime;comment:创建时间"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime;comment:更新时间"`
	Connect   *gorm.DB  `json:"-" gorm:"-"`
}
