package gglmm

import (
	"time"
)

// Model 基本模型类型
type Model struct {
	ID        int64      `json:"id" gorm:"primary_key;"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt"`
}
