package gglmm

import (
	"errors"
	"log"
	"time"

	"github.com/jinzhu/gorm"
)

var (
	// ErrParams --
	ErrParams = errors.New("参数错误")
	// ErrRecordNotFound --
	ErrRecordNotFound = errors.New("记录不存在")
	// ErrStoreFail --
	ErrStoreFail = errors.New("保存失败")
	// ErrStoreFailNotNewRecord --
	ErrStoreFailNotNewRecord = errors.New("保存失败，不是新记录")
	// ErrUpdateFail --
	ErrUpdateFail = errors.New("更新失败")
	// ErrUpdateFailID --
	ErrUpdateFailID = errors.New("更新失败，ID错误")
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

// ID --
func (gormDB *GormDB) ID(model interface{}) int64 {
	if entity, ok := model.(Entity); ok {
		return entity.UniqueID()
	}
	return 0
}

func (gormDB *GormDB) preloadDB(preloads []string) *gorm.DB {
	db := gormDB.DB
	for _, preload := range preloads {
		db = db.Preload(preload)
	}
	return db
}

// Get 单个查询
func (gormDB *GormDB) Get(model interface{}, request interface{}) error {
	if idRequest, ok := request.(*IDRequest); ok {
		return gormDB.getByID(model, idRequest)
	} else if idRequest, ok := request.(*IDRequest); ok {
		return gormDB.getByID(model, idRequest)
	} else if filterRequest, ok := request.(*FilterRequest); ok {
		return gormDB.getByFilter(model, filterRequest)
	} else if filterRequest, ok := request.(*FilterRequest); ok {
		return gormDB.getByFilter(model, filterRequest)
	} else if id, ok := request.(int64); ok {
		idRequest := &IDRequest{
			ID: id,
		}
		return gormDB.getByID(model, idRequest)
	}
	return ErrParams
}

// getByID 通过ID单个查询
func (gormDB *GormDB) getByID(model interface{}, idRequest *IDRequest) error {
	db := gormDB.preloadDB(idRequest.Preloads)
	if err := db.First(model, idRequest.ID).Error; err != nil {
		return err
	}
	if db.NewRecord(model) {
		return ErrRecordNotFound
	}
	return nil
}

// getByFilter 根据条件单个查询
func (gormDB *GormDB) getByFilter(model interface{}, filterRequest *FilterRequest) error {
	db := gormDB.preloadDB(filterRequest.Preloads)
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
func (gormDB *GormDB) List(slice interface{}, filterRequest *FilterRequest) error {
	db := gormDB.preloadDB(filterRequest.Preloads)
	db, err := gormFilterRequest(db, filterRequest)
	if err != nil {
		return err
	}
	if err = db.Find(slice).Error; err != nil {
		return err
	}
	return nil
}

// Page 根据条件分页查询
func (gormDB *GormDB) Page(pageResponse *PageResponse, pageRequest *PageRequest) error {
	db := gormDB.preloadDB(pageRequest.Preloads)
	db, err := gormFilterRequest(db, pageRequest.FilterRequest)
	if err != nil {
		return err
	}
	pagination := pageRequest.Pagination
	pageResponse.Pagination = pagination
	offset := (pagination.PageIndex - 1) * pagination.PageSize
	if err = db.Model(pageResponse.List).Count(&pageResponse.Pagination.Total).Limit(pageResponse.Pagination.PageSize).Offset(offset).Find(pageResponse.List).Error; err != nil {
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
	if gormDB.NewRecord(model) {
		return ErrStoreFail
	}
	return nil
}

// Update 更新整体
func (gormDB *GormDB) Update(model interface{}, id int64) error {
	modelID := gormDB.ID(model)
	if modelID != id {
		return ErrParams
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
func (gormDB *GormDB) UpdateFields(model interface{}, fields interface{}) error {
	modelID := gormDB.ID(model)
	if err := gormDB.Model(model).Updates(fields).Error; err != nil {
		return err
	}
	if err := gormDB.First(model, modelID).Error; err != nil {
		return err
	}
	return nil
}

// Remove 软删除
func (gormDB *GormDB) Remove(model interface{}, id int64) error {
	if err := gormDB.Delete(model, "id = ?", id).Error; err != nil {
		return err
	}
	if err := gormDB.Unscoped().First(model, id).Error; err != nil {
		return err
	}
	return nil
}

// Restore 恢复
func (gormDB *GormDB) Restore(model interface{}, id int64) error {
	if err := gormDB.Unscoped().Model(model).Where("id = ?", id).Update("deleted_at", nil).Error; err != nil {
		return err
	}
	if err := gormDB.First(model, id).Error; err != nil {
		return err
	}
	return nil
}

// Destroy 直接删除
func (gormDB *GormDB) Destroy(model interface{}, id int64) error {
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
