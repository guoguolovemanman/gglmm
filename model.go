package gglmm

import (
	"time"
)

// Model 基本模型类型
type Model struct {
	ID        int64      `json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt"`
}

// UniqueID --
func (model Model) UniqueID() int64 {
	return model.ID
}

// Entity --
type Entity interface {
	UniqueID() int64
}
