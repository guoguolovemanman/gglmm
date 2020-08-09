package gglmm

import (
	"errors"
	"log"
	"time"

	"github.com/jinzhu/gorm"
)

// error
var (
	ErrParams                = errors.New("参数错误")
	ErrStoreFailNotNewRecord = errors.New("保存失败，不是新记录")
	ErrUpdateID              = errors.New("更新失败")
)

// GormDB 服务
type GormDB struct {
	*gorm.DB
}

// NewGormDBConfig 新建服务
func NewGormDBConfig(config ConfigDB) *GormDB {
	return &GormDB{
		DB: GormOpenConfig(config),
	}
}

// NewGormDB --
func NewGormDB(dialect string, url string, maxOpen int, maxIdle int, connMaxLifetime time.Duration) *GormDB {
	return &GormDB{
		DB: GormOpen(dialect, url, maxOpen, maxIdle, connMaxLifetime),
	}
}

func (gormDB *GormDB) preloadsDB(preloads []string) *gorm.DB {
	db := gormDB.DB
	for _, preload := range preloads {
		db = db.Preload(preload)
	}
	return db
}

func (gormDB *GormDB) primaryKeyValue(model interface{}) uint64 {
	scope := gormDB.NewScope(model)
	key := scope.PrimaryKey()
	if key == "id" {
		value := scope.PrimaryKeyValue()
		if value, ok := value.(uint64); ok {
			return value
		}
	}
	return 0
}

// Get 单个查询
func (gormDB *GormDB) Get(model interface{}, request interface{}) error {
	switch request := request.(type) {
	case uint64:
		idRequest := &IDRequest{
			ID: request,
		}
		return gormDB.GetByID(model, idRequest)
	case IDRequest:
		return gormDB.GetByID(model, &request)
	case *IDRequest:
		return gormDB.GetByID(model, request)
	case FilterRequest:
		return gormDB.FirstByFilter(model, &request)
	case *FilterRequest:
		return gormDB.FirstByFilter(model, request)
	default:
		return ErrParams
	}
}

// GetByID 通过ID单个查询
func (gormDB *GormDB) GetByID(model interface{}, idRequest *IDRequest) error {
	db := gormDB.preloadsDB(idRequest.Preloads)
	if err := db.First(model, idRequest.ID).Error; err != nil {
		return err
	}
	return nil
}

// FirstByFilter 根据条件单个查询
func (gormDB *GormDB) FirstByFilter(model interface{}, filterRequest *FilterRequest) error {
	db := gormDB.preloadsDB(filterRequest.Preloads)
	db, err := gormFilterRequest(db, filterRequest)
	if err != nil {
		return err
	}
	if err = db.First(model).Error; err != nil {
		return err
	}
	return nil
}

// List 根据条件列表查询
func (gormDB *GormDB) List(models interface{}, filterRequest *FilterRequest) error {
	db := gormDB.preloadsDB(filterRequest.Preloads)
	db, err := gormFilterRequest(db, filterRequest)
	if err != nil {
		return err
	}
	if err = db.Find(models).Error; err != nil {
		return err
	}
	return nil
}

// Page 根据条件分页查询
func (gormDB *GormDB) Page(response *PageResponse, request *PageRequest) error {
	db := gormDB.preloadsDB(request.Preloads)
	db, err := gormFilterRequest(db, request.FilterRequest)
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
	if err = db.Model(response.List).Count(&response.Pagination.Total).Limit(response.Pagination.PageSize).Offset(offset).Find(response.List).Error; err != nil {
		return err
	}
	return nil
}

// Store 保存
func (gormDB *GormDB) Store(model interface{}) error {
	if !gormDB.NewRecord(model) {
		return ErrStoreFailNotNewRecord
	}
	if err := gormDB.Create(model).Error; err != nil {
		return err
	}
	return nil
}

// Update 更新整体
func (gormDB *GormDB) Update(model interface{}, id uint64) error {
	if id != gormDB.primaryKeyValue(model) {
		return ErrUpdateID
	}
	if err := gormDB.Save(model).Error; err != nil {
		return err
	}
	if err := gormDB.First(model, id).Error; err != nil {
		return err
	}
	return nil
}

// UpdateFields 更新多个属性
func (gormDB *GormDB) UpdateFields(model interface{}, id uint64, fields interface{}) error {
	if id != gormDB.primaryKeyValue(model) {
		return ErrUpdateID
	}
	if err := gormDB.Model(model).Updates(fields).Error; err != nil {
		return err
	}
	if err := gormDB.First(model, id).Error; err != nil {
		return err
	}
	return nil
}

// Remove 软删除
func (gormDB *GormDB) Remove(model interface{}, id uint64) error {
	if err := gormDB.Delete(model, "id = ?", id).Error; err != nil {
		return err
	}
	if err := gormDB.Unscoped().First(model, id).Error; err != nil {
		return err
	}
	return nil
}

// Restore 恢复
func (gormDB *GormDB) Restore(model interface{}, id uint64) error {
	if err := gormDB.Unscoped().Model(model).Where("id = ?", id).Update("deleted_at", nil).Error; err != nil {
		return err
	}
	if err := gormDB.First(model, id).Error; err != nil {
		return err
	}
	return nil
}

// Destroy 直接删除
func (gormDB *GormDB) Destroy(model interface{}, id uint64) error {
	return gormDB.Unscoped().Delete(model, "id = ?", id).Error
}

var gormDB *GormDB = nil

// RegisterGormDBConfig --
func RegisterGormDBConfig(config ConfigDB) {
	gormDB = NewGormDBConfig(config)
}

// RegisterGormDB --
func RegisterGormDB(dialect string, url string, maxOpen int, maxIdle int, connMaxLifetime time.Duration) {
	gormDB = NewGormDB(dialect, url, maxOpen, maxIdle, connMaxLifetime)
}

// CloseGormDB --
func CloseGormDB() {
	if gormDB != nil {
		gormDB.Close()
	}
}

// ErrGormDBNotRegister --
var ErrGormDBNotRegister = errors.New("请注册GromDB")

// DefaultGormDB --
func DefaultGormDB() *GormDB {
	if nil == gormDB {
		log.Fatal(ErrGormDBNotRegister)
	}
	return gormDB
}
