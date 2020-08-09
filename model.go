package gglmm

import "time"

// Model 基本模型类型
type Model struct {
	ID        uint64     `json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt"`
}
