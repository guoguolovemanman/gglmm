package gglmm

import (
	"errors"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

const (
	// ActionGetByID 根据ID拉取单个
	ActionGetByID = "GetByID"
	// ActionFirst 根据条件拉取单个
	ActionFirst = "First"
	// ActionList 列表
	ActionList = "List"
	// ActionPage 分页
	ActionPage = "page"
	// ActionStore 保存
	ActionStore = "Store"
	// ActionUpdate 更新整体
	ActionUpdate = "Update"
	// ActionUpdateFields 更新多个字段
	ActionUpdateFields = "UpdateFields"
	// ActionDestory 软删除
	ActionDestory = "Destory"
	// ActionRestore 恢复
	ActionRestore = "Resotre"
	// ActionRemove 直接删除
	ActionRemove = "Remove"
)

// IDRegexp ID正则表达式
const IDRegexp = "{id:[0-9]+}"

var (
	// ReadActions 读操作
	ReadActions = []string{ActionGetByID, ActionFirst, ActionList, ActionPage}
	// WriteActions 写操作
	WriteActions = []string{ActionStore, ActionUpdate, ActionUpdateFields}
	// DeleteActions 删除操作
	DeleteActions = []string{ActionDestory, ActionRestore}
	// AdminActions 管理操作
	AdminActions = []string{ActionList, ActionPage, ActionStore, ActionUpdate, ActionDestory, ActionRestore}
	// AllActions 全部操作
	AllActions = []string{ActionGetByID, ActionFirst, ActionList, ActionPage, ActionStore, ActionUpdate, ActionUpdateFields, ActionDestory, ActionRestore}
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

// Action --
func (service *HTTPService) Action(action string) (*HTTPAction, error) {
	var path string
	var handlerFunc http.HandlerFunc
	var method string
	switch action {
	case ActionGetByID:
		path = "/" + IDRegexp
		handlerFunc = service.GetByID
		method = "GET"
	case ActionFirst:
		path = "/first"
		handlerFunc = service.First
		method = "POST"
	case ActionList:
		path = "/list"
		handlerFunc = service.List
		method = "POST"
	case ActionPage:
		path = "/page"
		handlerFunc = service.Page
		method = "POST"
	case ActionStore:
		handlerFunc = service.Store
		method = "POST"
	case ActionUpdate:
		path = "/" + IDRegexp
		handlerFunc = service.Update
		method = "PUT"
	case ActionUpdateFields:
		path = "/" + IDRegexp + "/fields"
		handlerFunc = service.UpdateFields
		method = "PUT"
	case ActionDestory:
		path = "/" + IDRegexp
		handlerFunc = service.Destory
		method = "DELETE"
	case ActionRestore:
		path = "/" + IDRegexp
		handlerFunc = service.Restore
		method = "POST"
	case ActionRemove:
		path = "/" + IDRegexp + "/hard"
		handlerFunc = service.Remove
		method = "DELETE"
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
		RequestErrorResponse(err.Error()).JSON(w)
		return
	}
	model := reflect.New(service.modelType).Interface()
	if cacher != nil {
		if SupportCache(service.modelValue) {
			cacheKey := service.modelType.Name() + ":" + strconv.FormatInt(idRequest.ID, 10)
			if len(idRequest.Preloads) > 0 {
				cacheKey = cacheKey + ":" + strings.Join(idRequest.Preloads, "-")
			}
			if err := cacher.GetObj(cacheKey, model); err == nil {
				OkResponse().
					AddData(SingleKey(service.modelValue), model).
					JSON(w)
				return
			}
		}
	}
	if err = gormRepository.Get(model, idRequest); err != nil {
		FailResponse(err.Error()).JSON(w)
		return
	}
	if cacher != nil {
		if SupportCache(service.modelValue) {
			cacheKey := service.modelType.Name() + ":" + strconv.FormatInt(idRequest.ID, 10)
			if len(idRequest.Preloads) > 0 {
				cacheKey = cacheKey + ":" + strings.Join(idRequest.Preloads, "-")
			}
			cacher.Set(cacheKey, model)
		}
	}
	OkResponse().
		AddData(SingleKey(service.modelValue), model).
		JSON(w)
}

// First 单个
func (service *HTTPService) First(w http.ResponseWriter, r *http.Request) {
	filterRequest, err := DecodeFilterRequest(r)
	if err != nil {
		RequestErrorResponse(err.Error()).JSON(w)
		return
	}
	if service.filterFunc != nil {
		filterRequest.Filters = service.filterFunc(filterRequest.Filters, r)
	}
	model := reflect.New(service.modelType).Interface()
	if err = gormRepository.Get(model, filterRequest); err != nil {
		FailResponse(err.Error()).JSON(w)
		return
	}
	OkResponse().
		AddData(SingleKey(service.modelValue), model).
		JSON(w)
}

// List 列表
func (service *HTTPService) List(w http.ResponseWriter, r *http.Request) {
	filterRequest, err := DecodeFilterRequest(r)
	if err != nil {
		RequestErrorResponse(err.Error()).JSON(w)
		return
	}
	if service.filterFunc != nil {
		filterRequest.Filters = service.filterFunc(filterRequest.Filters, r)
	}
	list := reflect.New(reflect.SliceOf(service.modelType)).Interface()
	if err = gormRepository.List(list, filterRequest); err != nil {
		FailResponse(err.Error()).JSON(w)
		return
	}
	OkResponse().
		AddData(MultiKey(service.modelValue), list).
		JSON(w)
}

// Page 分页
func (service *HTTPService) Page(w http.ResponseWriter, r *http.Request) {
	pageRequest, err := DecodePageRequest(r)
	if err != nil {
		RequestErrorResponse(err.Error()).JSON(w)
		return
	}
	if service.filterFunc != nil {
		pageRequest.Filters = service.filterFunc(pageRequest.Filters, r)
	}
	pageResponse := PageResponse{}
	pageResponse.List = reflect.New(reflect.SliceOf(service.modelType)).Interface()
	if err = gormRepository.Page(&pageResponse, pageRequest); err != nil {
		FailResponse(err.Error()).JSON(w)
		return
	}
	OkResponse().
		AddData(MultiKey(service.modelValue), pageResponse.List).
		AddData("pagination", pageResponse.Pagination).
		JSON(w)
}

// Store 保存
func (service *HTTPService) Store(w http.ResponseWriter, r *http.Request) {
	model, err := DecodeModelPtr(r, service.modelType)
	if err != nil {
		RequestErrorResponse(err.Error()).JSON(w)
		return
	}
	if service.beforeStoreFunc != nil {
		model, err = service.beforeStoreFunc(model)
		if err != nil {
			FailResponse(err.Error()).JSON(w)
			return
		}
	}
	if err = gormRepository.Store(model); err != nil {
		FailResponse(err.Error()).JSON(w)
		return
	}
	OkResponse().
		AddData(SingleKey(service.modelValue), model).
		JSON(w)
}

// Update 更新整体
func (service *HTTPService) Update(w http.ResponseWriter, r *http.Request) {
	id, err := PathVarID(r)
	if err != nil {
		RequestErrorResponse(err.Error()).JSON(w)
		return
	}
	model, err := DecodeModelPtr(r, service.modelType)
	if err != nil {
		FailResponse(err.Error()).JSON(w)
		return
	}
	if service.beforeUpdateFunc != nil {
		model, id, err = service.beforeUpdateFunc(model, id)
		if err != nil {
			FailResponse(err.Error()).JSON(w)
			return
		}
	}
	if err = gormRepository.Update(model, id); err != nil {
		FailResponse(err.Error()).JSON(w)
		return
	}
	if cacher != nil {
		if SupportCache(service.modelValue) {
			cacheKey := service.modelType.Name() + ":" + strconv.FormatInt(id, 10)
			cacher.DelPattern(cacheKey)
		}
	}
	OkResponse().
		AddData(SingleKey(service.modelValue), model).
		JSON(w)
}

// UpdateFields 更新整体
func (service *HTTPService) UpdateFields(w http.ResponseWriter, r *http.Request) {
	id, err := PathVarID(r)
	if err != nil {
		RequestErrorResponse(err.Error()).JSON(w)
		return
	}
	model := reflect.New(service.modelType).Interface()
	if err := gormRepository.Get(model, id); err != nil {
		FailResponse(err.Error()).JSON(w)
		return
	}
	fields, err := DecodeModelPtr(r, service.modelType)
	if err != nil {
		FailResponse(err.Error()).JSON(w)
		return
	}
	if service.beforeUpdateFunc != nil {
		fields, id, err = service.beforeUpdateFunc(fields, id)
		if err != nil {
			FailResponse(err.Error()).JSON(w)
			return
		}
	}
	if err = gormRepository.UpdateFields(model, fields); err != nil {
		FailResponse(err.Error()).JSON(w)
		return
	}
	if cacher != nil {
		if SupportCache(service.modelValue) {
			cacheKey := service.modelType.Name() + ":" + strconv.FormatInt(id, 10)
			cacher.DelPattern(cacheKey)
		}
	}
	OkResponse().
		AddData(SingleKey(service.modelValue), model).
		JSON(w)
}

// Destory 软删除
func (service *HTTPService) Destory(w http.ResponseWriter, r *http.Request) {
	id, err := PathVarID(r)
	if err != nil {
		RequestErrorResponse(err.Error()).JSON(w)
		return
	}
	model := reflect.New(service.modelType).Interface()
	if err = gormRepository.Destroy(model, id); err != nil {
		FailResponse(err.Error()).JSON(w)
		return
	}
	if cacher != nil {
		if SupportCache(service.modelValue) {
			cacheKey := service.modelType.Name() + ":" + strconv.FormatInt(id, 10)
			cacher.DelPattern(cacheKey)
		}
	}
	OkResponse().
		AddData(SingleKey(service.modelValue), model).
		JSON(w)
}

// Restore 恢复
func (service *HTTPService) Restore(w http.ResponseWriter, r *http.Request) {
	id, err := PathVarID(r)
	if err != nil {
		RequestErrorResponse(err.Error()).JSON(w)
		return
	}
	model := reflect.New(service.modelType).Interface()
	if err = gormRepository.Restore(model, id); err != nil {
		FailResponse(err.Error()).JSON(w)
		return
	}
	if cacher != nil {
		if SupportCache(service.modelValue) {
			cacheKey := service.modelType.Name() + ":" + strconv.FormatInt(id, 10)
			cacher.DelPattern(cacheKey)
		}
	}
	OkResponse().
		AddData(SingleKey(service.modelValue), model).
		JSON(w)
}

// Remove 直接删除
func (service *HTTPService) Remove(w http.ResponseWriter, r *http.Request) {
	id, err := PathVarID(r)
	if err != nil {
		RequestErrorResponse(err.Error()).JSON(w)
		return
	}
	model := reflect.New(service.modelType).Interface()
	if err = gormRepository.Remove(model, id); err != nil {
		FailResponse(err.Error()).JSON(w)
		return
	}
	if cacher != nil {
		if SupportCache(service.modelValue) {
			cacheKey := service.modelType.Name() + ":" + strconv.FormatInt(id, 10)
			cacher.DelPattern(cacheKey)
		}
	}
	OkResponse().JSON(w)
}
