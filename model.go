package gglmm

import (
	"errors"
	"time"
)

// ErrModelType --
var ErrModelType = errors.New("模型类型错误")

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
