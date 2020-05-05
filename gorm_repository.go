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

// GormRepository 服务
type GormRepository struct {
	*gorm.DB
}

// NewGormRepositoryConfig 新建服务
func NewGormRepositoryConfig(config ConfigDB) *GormRepository {
	return &GormRepository{
		DB: NewGormDBConfig(config),
	}
}

// NewGormRepository --
func NewGormRepository(dialect string, url string, maxOpen int, maxIdle int, connMaxLifetime time.Duration) *GormRepository {
	return &GormRepository{
		DB: NewGormDB(dialect, url, maxOpen, maxIdle, connMaxLifetime),
	}
}

// ID --
func (repository *GormRepository) ID(model interface{}) int64 {
	if entity, ok := model.(Entity); ok {
		return entity.UniqueID()
	}
	return 0
}

func (repository *GormRepository) preloadDB(preloads []string) *gorm.DB {
	db := repository.DB
	for _, preload := range preloads {
		db = db.Preload(preload)
	}
	return db
}

// Get 单个查询
func (repository *GormRepository) Get(model interface{}, request interface{}) error {
	if idRequest, ok := request.(IDRequest); ok {
		return repository.getByID(model, idRequest)
	} else if idRequest, ok := request.(*IDRequest); ok {
		return repository.getByID(model, *idRequest)
	} else if filterRequest, ok := request.(FilterRequest); ok {
		return repository.getByFilter(model, filterRequest)
	} else if filterRequest, ok := request.(*FilterRequest); ok {
		return repository.getByFilter(model, *filterRequest)
	} else if id, ok := request.(int64); ok {
		idRequest := IDRequest{
			ID: id,
		}
		return repository.getByID(model, idRequest)
	}
	return ErrParams
}

// getByID 通过ID单个查询
func (repository *GormRepository) getByID(model interface{}, idRequest IDRequest) error {
	db := repository.preloadDB(idRequest.Preloads)
	if err := db.First(model, idRequest.ID).Error; err != nil {
		return err
	}
	if db.NewRecord(model) {
		return ErrRecordNotFound
	}
	return nil
}

// getByFilter 根据条件单个查询
func (repository *GormRepository) getByFilter(model interface{}, filterRequest FilterRequest) error {
	db := repository.preloadDB(filterRequest.Preloads)
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
func (repository *GormRepository) List(slice interface{}, filterRequest FilterRequest) error {
	db := repository.preloadDB(filterRequest.Preloads)
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
func (repository *GormRepository) Page(pageResponse *PageResponse, pageRequest PageRequest) error {
	db := repository.preloadDB(pageRequest.Preloads)
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
func (repository *GormRepository) Store(model interface{}) error {
	if !repository.NewRecord(model) {
		return ErrStoreFailNotNewRecord
	}
	if err := repository.Create(model).Error; err != nil {
		return err
	}
	if repository.NewRecord(model) {
		return ErrStoreFail
	}
	return nil
}

// Update 更新整体
func (repository *GormRepository) Update(model interface{}, id int64) error {
	modelID := repository.ID(model)
	if modelID != id {
		return ErrParams
	}
	if err := repository.Save(model).Error; err != nil {
		return err
	}
	if err := repository.First(model, id).Error; err != nil {
		return err
	}
	return nil
}

// UpdateFields 更新多个属性
func (repository *GormRepository) UpdateFields(model interface{}, fields interface{}) error {
	modelID := repository.ID(model)
	if err := repository.Model(model).Updates(fields).Error; err != nil {
		return err
	}
	if err := repository.First(model, modelID).Error; err != nil {
		return err
	}
	return nil
}

// Remove 软删除
func (repository *GormRepository) Remove(model interface{}, id int64) error {
	if err := repository.Delete(model, "id = ?", id).Error; err != nil {
		return err
	}
	if err := repository.Unscoped().First(model, id).Error; err != nil {
		return err
	}
	return nil
}

// Restore 恢复
func (repository *GormRepository) Restore(model interface{}, id int64) error {
	if err := repository.Unscoped().Model(model).Where("id = ?", id).Update("deleted_at", nil).Error; err != nil {
		return err
	}
	if err := repository.First(model, id).Error; err != nil {
		return err
	}
	return nil
}

// Destroy 直接删除
func (repository *GormRepository) Destroy(model interface{}, id int64) error {
	return repository.Unscoped().Delete(model, "id = ?", id).Error
}

var gormRepository *GormRepository = nil

// RegisterGormRepositoryConfig --
func RegisterGormRepositoryConfig(config ConfigDB) {
	gormRepository = NewGormRepositoryConfig(config)
}

// RegisterGormRepository --
func RegisterGormRepository(dialect string, url string, maxOpen int, maxIdle int, connMaxLifetime time.Duration) {
	gormRepository = NewGormRepository(dialect, url, maxOpen, maxIdle, connMaxLifetime)
}

// CloseGormRepository --
func CloseGormRepository() {
	if gormRepository != nil {
		gormRepository.Close()
	}
}

// ErrGormRepositoryNotRegister --
var ErrGormRepositoryNotRegister = errors.New("请注册GromRepository")

// DefaultGormRepository --
func DefaultGormRepository() *GormRepository {
	if nil == gormRepository {
		log.Fatal(ErrGormRepositoryNotRegister)
	}
	return gormRepository
}
