package gglmm

import (
	"errors"

	"github.com/jinzhu/gorm"
)

// error
var (
	ErrParams       = errors.New("参数错误")
	ErrNotNewRecord = errors.New("不是新记录")
	ErrUpdateID     = errors.New("更新失败")
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

// Begin --
func (gglmmDB *DB) Begin() *gorm.DB {
	return gglmmDB.gormDB.Begin()
}

func (gglmmDB *DB) preloadsGormDB(preloads []string) *gorm.DB {
	db := gglmmDB.gormDB
	for _, preload := range preloads {
		db = db.Preload(preload)
	}
	return db
}

func (gglmmDB *DB) primaryKeyValue(model interface{}) uint64 {
	scope := gglmmDB.gormDB.NewScope(model)
	key := scope.PrimaryKey()
	if key == "id" {
		value := scope.PrimaryKeyValue()
		if value, ok := value.(uint64); ok {
			return value
		}
	}
	return 0
}

// Create 保存
func (gglmmDB *DB) Create(model interface{}) error {
	if !gglmmDB.gormDB.NewRecord(model) {
		return ErrNotNewRecord
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
		return ErrParams
	}
}

func (gglmmDB *DB) first(model interface{}, idRequest *IDRequest) error {
	gormDB := gglmmDB.preloadsGormDB(idRequest.Preloads)
	if err := gormDB.First(model, idRequest.ID).Error; err != nil {
		return err
	}
	return nil
}

func (gglmmDB *DB) firstByFilter(model interface{}, filterRequest *FilterRequest) error {
	gormDB := gglmmDB.preloadsGormDB(filterRequest.Preloads)
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
	gormDB := gglmmDB.preloadsGormDB(filterRequest.Preloads)
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
	gormDB := gglmmDB.preloadsGormDB(request.Preloads)
	gormDB, err := gormFilterRequest(gormDB, request.FilterRequest)
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
func (gglmmDB *DB) Update(model interface{}, id uint64) error {
	if id != gglmmDB.primaryKeyValue(model) {
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
func (gglmmDB *DB) Updates(model interface{}, id uint64, fields map[string]interface{}) error {
	if id != gglmmDB.primaryKeyValue(model) {
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
func (gglmmDB *DB) Remove(model interface{}, id uint64) error {
	if err := gglmmDB.gormDB.Delete(model, "id = ?", id).Error; err != nil {
		return err
	}
	if err := gglmmDB.gormDB.Unscoped().First(model, id).Error; err != nil {
		return err
	}
	return nil
}

// Restore 恢复
func (gglmmDB *DB) Restore(model interface{}, id uint64) error {
	if err := gglmmDB.gormDB.Unscoped().Model(model).Where("id = ?", id).Update("deleted_at", nil).Error; err != nil {
		return err
	}
	if err := gglmmDB.gormDB.First(model, id).Error; err != nil {
		return err
	}
	return nil
}

// Destroy 直接删除
func (gglmmDB *DB) Destroy(model interface{}, id uint64) error {
	if err := gglmmDB.gormDB.Unscoped().Delete(model, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}
