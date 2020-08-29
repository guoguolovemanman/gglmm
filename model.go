package gglmm

import "time"

// Entity --
type Entity interface {
	PrimaryKeyValue() uint64
	SetPrimaryKeyValue(uint64)
}

// Model 基本模型类型
type Model struct {
	ID        uint64     `json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt"`
}

// PrimaryKeyValue --
func (model *Model) PrimaryKeyValue() uint64 {
	return model.ID
}

// SetPrimaryKeyValue --
func (model *Model) SetPrimaryKeyValue(id uint64) {
	model.ID = id
}

// PrimaryKeyValue --
func PrimaryKeyValue(model interface{}) uint64 {
	if entity, ok := model.(Entity); ok {
		return entity.PrimaryKeyValue()
	}
	return 0
}

// SetPrimaryKeyValue --
func SetPrimaryKeyValue(model interface{}, id uint64) {
	if entity, ok := model.(Entity); ok {
		entity.SetPrimaryKeyValue(id)
	}
}
