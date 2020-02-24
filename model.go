package gglmm

import (
	"time"
)

// Model 基本模型类型
type Model struct {
	ID        int64      `gorm:"primary_key;" json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt"`
	Status    int8       `gorm:"not null;" json:"status"`
}
