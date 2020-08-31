package gglmm

import (
	"errors"

	"github.com/jinzhu/gorm"
)

// error
var (
	ErrCreateNotNewRecord = errors.New("新建失败，已存在主键")
	ErrUpdateID           = errors.New("更新失败，请设置主键")
	ErrDeleteID           = errors.New("删除失败，请设置主键")
)

// DB --
type DB struct {
	gormDB *gorm.DB
}

// NewDB 新建DB
func NewDB() *DB {
	return &DB{
		gormDB: DefaultGormDB(),
	}
}

// GormDB --
func (gglmmDB *DB) GormDB() *gorm.DB {
	return gglmmDB.gormDB
}

// NewRecord --
func (gglmmDB *DB) NewRecord(model interface{}) bool {
	return gglmmDB.gormDB.NewRecord(model)
}

func (gglmmDB *DB) primaryKeyValue(model interface{}) uint64 {
	if gglmmDB.gormDB.NewRecord(model) {
		return 0
	}
	return PrimaryKeyValue(model)
}

// Begin --
func (gglmmDB *DB) Begin() *gorm.DB {
	return gglmmDB.gormDB.Begin()
}

// Create 新建
func (gglmmDB *DB) Create(model interface{}) error {
	if !gglmmDB.gormDB.NewRecord(model) {
		return ErrCreateNotNewRecord
	}
	if err := gglmmDB.gormDB.Create(model).Error; err != nil {
		return err
	}
	return nil
}

// First 查询
func (gglmmDB *DB) First(model interface{}, request interface{}) error {
	switch request := request.(type) {
	case uint64:
		idRequest := IDRequest{
			ID: request,
		}
		return gglmmDB.first(model, &idRequest)
	case IDRequest:
		return gglmmDB.first(model, &request)
	case *IDRequest:
		return gglmmDB.first(model, request)
	case FilterRequest:
		return gglmmDB.firstByFilter(model, &request)
	case *FilterRequest:
		return gglmmDB.firstByFilter(model, request)
	default:
		return ErrParameter
	}
}

func (gglmmDB *DB) first(model interface{}, idRequest *IDRequest) error {
	gormDB := gormPreloads(gglmmDB.gormDB, idRequest.Preloads)
	if err := gormDB.First(model, idRequest.ID).Error; err != nil {
		return err
	}
	return nil
}

func (gglmmDB *DB) firstByFilter(model interface{}, filterRequest *FilterRequest) error {
	gormDB := gormPreloads(gglmmDB.gormDB, filterRequest.Preloads)
	gormDB, err := gormFilterRequest(gormDB, filterRequest)
	if err != nil {
		return err
	}
	if err := gormDB.First(model).Error; err != nil {
		return err
	}
	return nil
}

// List 根据条件列表查询
func (gglmmDB *DB) List(models interface{}, filterRequest *FilterRequest) error {
	gormDB := gormPreloads(gglmmDB.gormDB, filterRequest.Preloads)
	gormDB, err := gormFilterRequest(gormDB, filterRequest)
	if err != nil {
		return err
	}
	if err := gormDB.Find(models).Error; err != nil {
		return err
	}
	return nil
}

// Page 根据条件分页查询
func (gglmmDB *DB) Page(response *PageResponse, request *PageRequest) error {
	gormDB := gormPreloads(gglmmDB.gormDB, request.Preloads)
	gormDB, err := gormFilterRequest(gormDB, &request.FilterRequest)
	if err != nil {
		return err
	}
	pageIndex := request.Pagination.PageIndex
	if pageIndex == 0 {
		pageIndex = FirstPageIndex
	}
	pageSize := request.Pagination.PageSize
	if pageSize == 0 {
		pageSize = DefaultPageSize
	}
	response.Pagination.PageIndex = pageIndex
	response.Pagination.PageSize = pageSize
	offset := (pageIndex - 1) * pageSize
	if err = gormDB.Model(response.List).Count(&response.Pagination.Total).Limit(response.Pagination.PageSize).Offset(offset).Find(response.List).Error; err != nil {
		return err
	}
	return nil
}

// Update 更新整体
func (gglmmDB *DB) Update(model interface{}) error {
	id := gglmmDB.primaryKeyValue(model)
	if id <= 0 {
		return ErrUpdateID
	}
	if err := gglmmDB.gormDB.Save(model).Error; err != nil {
		return err
	}
	if err := gglmmDB.gormDB.First(model, id).Error; err != nil {
		return err
	}
	return nil
}

// Updates 更新多个属性
func (gglmmDB *DB) Updates(model interface{}, fields map[string]interface{}) error {
	id := gglmmDB.primaryKeyValue(model)
	if id <= 0 {
		return ErrUpdateID
	}
	if err := gglmmDB.gormDB.Model(model).Updates(fields).Error; err != nil {
		return err
	}
	if err := gglmmDB.gormDB.First(model, id).Error; err != nil {
		return err
	}
	return nil
}

// Remove 软删除
func (gglmmDB *DB) Remove(model interface{}) error {
	id := gglmmDB.primaryKeyValue(model)
	if id <= 0 {
		return ErrDeleteID
	}
	if err := gglmmDB.gormDB.Delete(model).Error; err != nil {
		return err
	}
	if err := gglmmDB.gormDB.Unscoped().First(model, id).Error; err != nil {
		return err
	}
	return nil
}

// Restore 恢复
func (gglmmDB *DB) Restore(model interface{}) error {
	id := gglmmDB.primaryKeyValue(model)
	if id <= 0 {
		return ErrDeleteID
	}
	if err := gglmmDB.gormDB.Unscoped().Model(model).Update("deleted_at", nil).Error; err != nil {
		return err
	}
	if err := gglmmDB.gormDB.First(model, id).Error; err != nil {
		return err
	}
	return nil
}

// Destroy 直接删除
func (gglmmDB *DB) Destroy(model interface{}) error {
	id := gglmmDB.primaryKeyValue(model)
	if id <= 0 {
		return ErrDeleteID
	}
	if err := gglmmDB.gormDB.Unscoped().Delete(model).Error; err != nil {
		return err
	}
	return nil
}
