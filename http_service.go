package gglmm

import (
	"errors"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

// RESTAction 操作
type RESTAction uint8

const (
	// RESTActionGetByID 根据ID拉取单个
	RESTActionGetByID RESTAction = iota
	// RESTActionFirst 根据条件拉取单个
	RESTActionFirst
	// RESTActionList 列表
	RESTActionList
	// RESTActionPage 分页
	RESTActionPage
	// RESTActionStore 保存
	RESTActionStore
	// RESTActionUpdate 更新整体
	RESTActionUpdate
	// RESTActionUpdateFields 更新多个字段
	RESTActionUpdateFields
	// RESTActionDestory 软删除
	RESTActionDestory
	// RESTActionRestore 恢复
	RESTActionRestore
)

// IDRegexp ID正则表达式
const IDRegexp = "{id:[0-9]+}"

var (
	// RESTRead 读操作
	RESTRead = []RESTAction{RESTActionGetByID, RESTActionFirst, RESTActionList, RESTActionPage}
	// RESTWrite 写操作
	RESTWrite = []RESTAction{RESTActionStore, RESTActionUpdate, RESTActionUpdateFields}
	// RESTDelete 删除操作
	RESTDelete = []RESTAction{RESTActionDestory, RESTActionRestore}
	// RESTAdmin 管理操作
	RESTAdmin = []RESTAction{RESTActionList, RESTActionPage, RESTActionStore, RESTActionUpdate, RESTActionDestory, RESTActionRestore}
	// RESTAll 全部操作
	RESTAll = []RESTAction{RESTActionGetByID, RESTActionFirst, RESTActionList, RESTActionPage, RESTActionStore, RESTActionUpdate, RESTActionUpdateFields, RESTActionDestory, RESTActionRestore}
	// ErrAction --
	ErrAction = errors.New("不支持Action")
)

// FilterFunc 过滤函数
type FilterFunc func(filters []Filter, r *http.Request) []Filter

// BeforeStoreFunc 保存前调用
type BeforeStoreFunc func(model interface{}) (interface{}, error)

// BeforeUpdateFunc 更新前调用
type BeforeUpdateFunc func(model interface{}, id int64) (interface{}, int64, error)

// HTTPService HTTP服务
type HTTPService struct {
	modelType        reflect.Type
	modelValue       reflect.Value
	filterFunc       FilterFunc
	beforeStoreFunc  BeforeStoreFunc
	beforeUpdateFunc BeforeUpdateFunc
}

// NewHTTPService 新建HTTP服务
func NewHTTPService(model interface{}) *HTTPService {
	if gormRepository == nil {
		log.Fatal(ErrGormRepositoryNotRegister)
	}
	return &HTTPService{
		modelType:  reflect.TypeOf(model),
		modelValue: reflect.ValueOf(model),
	}
}

// HandleFilterFunc 设置过滤参数函数
func (service *HTTPService) HandleFilterFunc(handler FilterFunc) {
	service.filterFunc = handler
}

// HandleBeforeStoreFunc 设置保存前执行函数
func (service *HTTPService) HandleBeforeStoreFunc(handler BeforeStoreFunc) {
	service.beforeStoreFunc = handler
}

// HandleBeforeUpdateFunc 设置更新前执行函数
func (service *HTTPService) HandleBeforeUpdateFunc(handler BeforeUpdateFunc) {
	service.beforeUpdateFunc = handler
}

// CustomActions --
func (service *HTTPService) CustomActions() ([]*HTTPAction, error) {
	return nil, nil
}

// RESTAction --
func (service *HTTPService) RESTAction(restAction RESTAction) (*HTTPAction, error) {
	var path string
	var handlerFunc http.HandlerFunc
	var method string
	switch restAction {
	case RESTActionGetByID:
		path = "/" + IDRegexp
		handlerFunc = service.GetByID
		method = "GET"
	case RESTActionFirst:
		path = "/first"
		handlerFunc = service.First
		method = "POST"
	case RESTActionList:
		path = "/list"
		handlerFunc = service.List
		method = "POST"
	case RESTActionPage:
		path = "/page"
		handlerFunc = service.Page
		method = "POST"
	case RESTActionStore:
		handlerFunc = service.Store
		method = "POST"
	case RESTActionUpdate:
		path = "/" + IDRegexp
		handlerFunc = service.Update
		method = "PUT"
	case RESTActionUpdateFields:
		path = "/" + IDRegexp + "/fields"
		handlerFunc = service.UpdateFields
		method = "PUT"
	case RESTActionDestory:
		path = "/" + IDRegexp
		handlerFunc = service.Destory
		method = "DELETE"
	case RESTActionRestore:
		path = "/" + IDRegexp
		handlerFunc = service.Restore
		method = "POST"
	}
	if handlerFunc != nil {
		return NewHTTPAction(path, handlerFunc, method), nil
	}
	return nil, ErrAction
}

// GetByID 单个
func (service *HTTPService) GetByID(w http.ResponseWriter, r *http.Request) {
	idRequest, err := DecodeIDRequest(r)
	if err != nil {
		NewFailResponse(err.Error()).WriteJSON(w)
		return
	}
	model := reflect.New(service.modelType).Interface()
	if cacher != nil {
		if ReflectCache(service.modelValue) {
			cacheKey := service.modelType.Name() + ":" + strconv.FormatInt(idRequest.ID, 10)
			if len(idRequest.Preloads) > 0 {
				cacheKey = cacheKey + ":" + strings.Join(idRequest.Preloads, "-")
			}
			if err := cacher.GetObj(cacheKey, model); err == nil {
				NewSuccessResponse().
					AddData(ReflectSingleKey(service.modelValue), model).
					WriteJSON(w)
				return
			}
		}
	}
	if err = gormRepository.Get(model, idRequest); err != nil {
		NewFailResponse(err.Error()).WriteJSON(w)
		return
	}
	if cacher != nil {
		if ReflectCache(service.modelValue) {
			cacheKey := service.modelType.Name() + ":" + strconv.FormatInt(idRequest.ID, 10)
			if len(idRequest.Preloads) > 0 {
				cacheKey = cacheKey + ":" + strings.Join(idRequest.Preloads, "-")
			}
			cacher.Set(cacheKey, model)
		}
	}
	NewSuccessResponse().
		AddData(ReflectSingleKey(service.modelValue), model).
		WriteJSON(w)
}

// First 单个
func (service *HTTPService) First(w http.ResponseWriter, r *http.Request) {
	filterRequest, err := DecodeFilterRequest(r)
	if err != nil {
		NewFailResponse(err.Error()).WriteJSON(w)
		return
	}
	if service.filterFunc != nil {
		filterRequest.Filters = service.filterFunc(filterRequest.Filters, r)
	}
	model := reflect.New(service.modelType).Interface()
	if err = gormRepository.Get(model, filterRequest); err != nil {
		NewFailResponse(err.Error()).WriteJSON(w)
		return
	}
	NewSuccessResponse().
		AddData(ReflectSingleKey(service.modelValue), model).
		WriteJSON(w)
}

// List 列表
func (service *HTTPService) List(w http.ResponseWriter, r *http.Request) {
	filterRequest, err := DecodeFilterRequest(r)
	if err != nil {
		NewFailResponse(err.Error()).WriteJSON(w)
		return
	}
	if service.filterFunc != nil {
		filterRequest.Filters = service.filterFunc(filterRequest.Filters, r)
	}
	list := reflect.New(reflect.SliceOf(service.modelType)).Interface()
	if err = gormRepository.List(list, filterRequest); err != nil {
		NewFailResponse(err.Error()).WriteJSON(w)
		return
	}
	NewSuccessResponse().
		AddData(ReflectMultiKey(service.modelValue), list).
		WriteJSON(w)
}

// Page 分页
func (service *HTTPService) Page(w http.ResponseWriter, r *http.Request) {
	pageRequest, err := DecodePageRequest(r)
	if err != nil {
		NewFailResponse(err.Error()).WriteJSON(w)
		return
	}
	if service.filterFunc != nil {
		pageRequest.Filters = service.filterFunc(pageRequest.Filters, r)
	}
	pageResponse := PageResponse{}
	pageResponse.List = reflect.New(reflect.SliceOf(service.modelType)).Interface()
	if err = gormRepository.Page(&pageResponse, pageRequest); err != nil {
		NewFailResponse(err.Error()).WriteJSON(w)
		return
	}
	NewSuccessResponse().
		AddData(ReflectMultiKey(service.modelValue), pageResponse.List).
		AddData("pagination", pageResponse.Pagination).
		WriteJSON(w)
}

// Store 保存
func (service *HTTPService) Store(w http.ResponseWriter, r *http.Request) {
	model, err := DecodeModelPtr(r, service.modelType)
	if err != nil {
		NewFailResponse(err.Error()).WriteJSON(w)
		return
	}
	if service.beforeStoreFunc != nil {
		model, err = service.beforeStoreFunc(model)
		if err != nil {
			NewFailResponse(err.Error()).WriteJSON(w)
			return
		}
	}
	if err = gormRepository.Store(model); err != nil {
		NewFailResponse(err.Error()).WriteJSON(w)
		return
	}
	NewSuccessResponse().
		AddData(ReflectSingleKey(service.modelValue), model).
		WriteJSON(w)
}

// Update 更新整体
func (service *HTTPService) Update(w http.ResponseWriter, r *http.Request) {
	id, err := MuxVarID(r)
	if err != nil {
		NewFailResponse(err.Error()).WriteJSON(w)
		return
	}
	model, err := DecodeModelPtr(r, service.modelType)
	if err != nil {
		NewFailResponse(err.Error()).WriteJSON(w)
		return
	}
	if service.beforeUpdateFunc != nil {
		model, id, err = service.beforeUpdateFunc(model, id)
		if err != nil {
			NewFailResponse(err.Error()).WriteJSON(w)
			return
		}
	}
	if err = gormRepository.Update(model, id); err != nil {
		NewFailResponse(err.Error()).WriteJSON(w)
		return
	}
	if cacher != nil {
		if ReflectCache(service.modelValue) {
			cacheKey := service.modelType.Name() + ":" + strconv.FormatInt(id, 10)
			cacher.DelPattern(cacheKey)
		}
	}
	NewSuccessResponse().
		AddData(ReflectSingleKey(service.modelValue), model).
		WriteJSON(w)
}

// UpdateFields 更新整体
func (service *HTTPService) UpdateFields(w http.ResponseWriter, r *http.Request) {
	id, err := MuxVarID(r)
	if err != nil {
		NewFailResponse(err.Error()).WriteJSON(w)
		return
	}
	model := reflect.New(service.modelType).Interface()
	if err := gormRepository.Get(model, id); err != nil {
		NewFailResponse(err.Error()).WriteJSON(w)
		return
	}
	fields, err := DecodeModelPtr(r, service.modelType)
	if err != nil {
		NewFailResponse(err.Error()).WriteJSON(w)
		return
	}
	if service.beforeUpdateFunc != nil {
		fields, id, err = service.beforeUpdateFunc(fields, id)
		if err != nil {
			NewFailResponse(err.Error()).WriteJSON(w)
			return
		}
	}
	if err = gormRepository.UpdateFields(model, fields); err != nil {
		NewFailResponse(err.Error()).WriteJSON(w)
		return
	}
	if cacher != nil {
		if ReflectCache(service.modelValue) {
			cacheKey := service.modelType.Name() + ":" + strconv.FormatInt(id, 10)
			cacher.DelPattern(cacheKey)
		}
	}
	NewSuccessResponse().
		AddData(ReflectSingleKey(service.modelValue), model).
		WriteJSON(w)
}

// Destory 删除
func (service *HTTPService) Destory(w http.ResponseWriter, r *http.Request) {
	id, err := MuxVarID(r)
	if err != nil {
		NewFailResponse(err.Error()).WriteJSON(w)
		return
	}
	model := reflect.New(service.modelType).Interface()
	if err = gormRepository.Destroy(model, id); err != nil {
		NewFailResponse(err.Error()).WriteJSON(w)
		return
	}
	if cacher != nil {
		if ReflectCache(service.modelValue) {
			cacheKey := service.modelType.Name() + ":" + strconv.FormatInt(id, 10)
			cacher.DelPattern(cacheKey)
		}
	}
	NewSuccessResponse().
		AddData(ReflectSingleKey(service.modelValue), model).
		WriteJSON(w)
}

// Restore 恢复
func (service *HTTPService) Restore(w http.ResponseWriter, r *http.Request) {
	id, err := MuxVarID(r)
	if err != nil {
		NewFailResponse(err.Error()).WriteJSON(w)
		return
	}
	model := reflect.New(service.modelType).Interface()
	if err = gormRepository.Restore(model, id); err != nil {
		NewFailResponse(err.Error()).WriteJSON(w)
		return
	}
	if cacher != nil {
		if ReflectCache(service.modelValue) {
			cacheKey := service.modelType.Name() + ":" + strconv.FormatInt(id, 10)
			cacher.DelPattern(cacheKey)
		}
	}
	NewSuccessResponse().
		AddData(ReflectSingleKey(service.modelValue), model).
		WriteJSON(w)
}
