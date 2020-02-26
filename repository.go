package gglmm

import (
	"errors"
	"reflect"
	"strconv"

	"github.com/jinzhu/gorm"
)

var (
	// ErrRecordNotFound --
	ErrRecordNotFound = errors.New("记录不存在")
	// ErrSaveFail --
	ErrSaveFail = errors.New("保存失败")
	// ErrSaveFailNotNewRecord --
	ErrSaveFailNotNewRecord = errors.New("保存失败，不是新记录")
	// ErrUpdateFail --
	ErrUpdateFail = errors.New("更新失败")
	// ErrUpdateFailNotNewRecord --
	ErrUpdateFailNotNewRecord = errors.New("更新失败，不是新记录")
)

// Repository 服务
type Repository struct {
	modelType  reflect.Type
	modelValue reflect.Value
}

// NewRepository 新建服务
func NewRepository(model interface{}) *Repository {
	return &Repository{
		modelType:  reflect.TypeOf(model),
		modelValue: reflect.ValueOf(model),
	}
}

// Begin 开始事务
func (repository *Repository) Begin() *gorm.DB {
	return gormDB.Begin()
}

func (repository *Repository) preloadDB(preloads []string) *gorm.DB {
	db := gormDB
	for _, preload := range preloads {
		db = db.Preload(preload)
	}
	return db
}

// Get 通过ID单个查询
func (repository *Repository) Get(id int64, model interface{}, preloads []string) error {
	if cacher != nil {
		if reflectCache(repository.modelValue) {
			cacheKey := repository.modelType.Name() + "#" + strconv.FormatInt(id, 10)
			if err := cacher.GetObj(cacheKey, model); err == nil {
				return nil
			}
		}
	}

	db := repository.preloadDB(preloads)
	if err := db.First(model, id).Error; err != nil {
		return err
	}
	if db.NewRecord(model) {
		return ErrRecordNotFound
	}

	if cacher != nil {
		if reflectCache(repository.modelValue) {
			cacheKey := repository.modelType.Name() + "#" + strconv.FormatInt(id, 10)
			cacher.Set(cacheKey, model)
		}
	}

	return nil
}

// First 根据条件单个查询
func (repository *Repository) First(filterRequest FilterRequest, model interface{}) error {
	db := repository.preloadDB(filterRequest.Preloads)
	db, err := gormSetupFilterRequest(db, filterRequest)
	if err != nil {
		return err
	}
	if err = db.First(model).Error; err != nil {
		return err
	}
	return nil
}

// List 根据条件列表查询
func (repository *Repository) List(filterRequest FilterRequest, slice interface{}) error {
	db := repository.preloadDB(filterRequest.Preloads)
	db, err := gormSetupFilterRequest(db, filterRequest)
	if err != nil {
		return err
	}
	if err = db.Find(slice).Error; err != nil {
		return err
	}
	return nil
}

// Page 根据条件分页查询
func (repository *Repository) Page(pageRequest PageRequest, pageResponse *PageResponse) error {
	db := repository.preloadDB(pageRequest.Preloads)
	db, err := gormSetupFilterRequest(db, pageRequest.FilterRequest)
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
func (repository *Repository) Store(model interface{}) error {
	if !gormDB.NewRecord(model) {
		return ErrSaveFailNotNewRecord
	}
	if err := gormDB.Create(model).Error; err != nil {
		return err
	}
	if gormDB.NewRecord(model) {
		return ErrSaveFail
	}
	return nil
}

// Update 更新
func (repository *Repository) Update(id int64, model interface{}) error {
	record := ReflectNew(repository.modelType)
	if err := repository.Get(id, record, nil); err != nil {
		return err
	}

	if err := gormDB.Model(record).Updates(model).Error; err != nil {
		return err
	}

	if cacher == nil {
		if reflectCache(repository.modelValue) {
			cacheKey := repository.modelType.Name() + "#" + strconv.FormatInt(id, 10)
			cacher.Del(cacheKey)
		}
	}

	if err := repository.Get(id, model, nil); err != nil {
		return err
	}
	return nil
}

// Destroy 删除
func (repository *Repository) Destroy(id int64, model interface{}) error {
	if err := gormDB.Delete(model, "id = ?", id).Error; err != nil {
		return err
	}
	if err := gormDB.Unscoped().First(model, id).Error; err != nil {
		return err
	}
	return nil
}

// Restore 恢复
func (repository *Repository) Restore(id int64, model interface{}) error {
	if err := gormDB.Unscoped().Model(model).Where("id = ?", id).Update("deleted_at", nil).Error; err != nil {
		return err
	}
	if err := gormDB.First(model, id).Error; err != nil {
		return err
	}
	return nil
}
