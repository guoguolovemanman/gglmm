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
	// RESTActionUpdate 更新
	RESTActionUpdate
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
	RESTWrite = []RESTAction{RESTActionStore, RESTActionUpdate}
	// RESTDelete 删除操作
	RESTDelete = []RESTAction{RESTActionDestory, RESTActionRestore}
	// RESTAdmin 管理操作
	RESTAdmin = []RESTAction{RESTActionList, RESTActionPage, RESTActionStore, RESTActionUpdate, RESTActionDestory, RESTActionRestore}
	// RESTAll 全部操作
	RESTAll = []RESTAction{RESTActionGetByID, RESTActionFirst, RESTActionList, RESTActionPage, RESTActionStore, RESTActionUpdate, RESTActionDestory, RESTActionRestore}
	// ErrAction --
	ErrAction = errors.New("不支持Action")
)

// FilterFunc 过滤函数
type FilterFunc func(filters []Filter, r *http.Request) []Filter

// BeforeStoreFunc 保存前调用
type BeforeStoreFunc func(model interface{}) interface{}

// BeforeUpdateFunc 更新前调用
type BeforeUpdateFunc func(model interface{}, id int64) (interface{}, int64)

// RESTHTTPService HTTP服务
type RESTHTTPService struct {
	modelType        reflect.Type
	modelValue       reflect.Value
	filterFunc       FilterFunc
	beforeStoreFunc  BeforeStoreFunc
	beforeUpdateFunc BeforeUpdateFunc
}

// NewRESTHTTPService 新建HTTP服务
func NewRESTHTTPService(model interface{}) *RESTHTTPService {
	if gormRepository == nil {
		log.Fatal(ErrGormRepositoryNotRegister)
	}
	return &RESTHTTPService{
		modelType:  reflect.TypeOf(model),
		modelValue: reflect.ValueOf(model),
	}
}

// HandleFilterFunc 设置过滤参数函数
func (service *RESTHTTPService) HandleFilterFunc(handler FilterFunc) {
	service.filterFunc = handler
}

// HandleBeforeStoreFunc 设置保存前执行函数
func (service *RESTHTTPService) HandleBeforeStoreFunc(handler BeforeStoreFunc) {
	service.beforeStoreFunc = handler
}

// HandleBeforeUpdateFunc 设置更新前执行函数
func (service *RESTHTTPService) HandleBeforeUpdateFunc(handler BeforeUpdateFunc) {
	service.beforeUpdateFunc = handler
}

// CustomActions --
func (service *RESTHTTPService) CustomActions() ([]*HTTPAction, error) {
	return nil, nil
}

// RESTAction --
func (service *RESTHTTPService) RESTAction(restAction RESTAction) (*HTTPAction, error) {
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
func (service *RESTHTTPService) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := MuxVarID(r)
	if err != nil {
		NewFailResponse(err.Error()).WriteJSON(w)
		return
	}
	model := reflect.New(service.modelType).Interface()
	preloads := parseQueryPreloads(r)
	if cacher != nil {
		if ReflectCache(service.modelValue) {
			cacheKey := service.modelType.Name() + ":" + strconv.FormatInt(id, 10)
			if len(preloads) > 0 {
				cacheKey = cacheKey + ":" + strings.Join(preloads, "-")
			}
			if err := cacher.GetObj(cacheKey, model); err == nil {
				NewSuccessResponse().
					AddData(ReflectSingleKey(service.modelValue), model).
					WriteJSON(w)
				return
			}
		}
	}
	idReuest := IDRequest{
		ID:       id,
		Preloads: preloads,
	}
	if err = gormRepository.Get(model, idReuest); err != nil {
		NewFailResponse(err.Error()).WriteJSON(w)
		return
	}
	if cacher != nil {
		if ReflectCache(service.modelValue) {
			cacheKey := service.modelType.Name() + ":" + strconv.FormatInt(id, 10)
			if len(preloads) > 0 {
				cacheKey = cacheKey + ":" + strings.Join(preloads, "-")
			}
			cacher.Set(cacheKey, model)
		}
	}
	NewSuccessResponse().
		AddData(ReflectSingleKey(service.modelValue), model).
		WriteJSON(w)
}

// First 单个
func (service *RESTHTTPService) First(w http.ResponseWriter, r *http.Request) {
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
func (service *RESTHTTPService) List(w http.ResponseWriter, r *http.Request) {
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
func (service *RESTHTTPService) Page(w http.ResponseWriter, r *http.Request) {
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
func (service *RESTHTTPService) Store(w http.ResponseWriter, r *http.Request) {
	model, err := DecodeModelPtr(r, service.modelType)
	if err != nil {
		NewFailResponse(err.Error()).WriteJSON(w)
		return
	}
	if service.beforeStoreFunc != nil {
		model = service.beforeStoreFunc(model)
		if model == nil {
			NewFailResponse("重设失败").WriteJSON(w)
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

// Update 更新
func (service *RESTHTTPService) Update(w http.ResponseWriter, r *http.Request) {
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
		model, id = service.beforeUpdateFunc(model, id)
		if model == nil || id == 0 {
			NewFailResponse("重设失败").WriteJSON(w)
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

// Destory 删除
func (service *RESTHTTPService) Destory(w http.ResponseWriter, r *http.Request) {
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
func (service *RESTHTTPService) Restore(w http.ResponseWriter, r *http.Request) {
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
