package gglmm

import (
	"errors"
	"net/http"

	"github.com/jinzhu/gorm"
)

// RESTAction 操作
type RESTAction uint8

const (
	// RESTActionGetByID 根据ID拉取单个
	RESTActionGetByID RESTAction = iota
	// RESTActionGet 根据条件拉取单个
	RESTActionGet
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
	RESTRead = []RESTAction{RESTActionGetByID, RESTActionGet, RESTActionList, RESTActionPage}
	// RESTWrite 写操作
	RESTWrite = []RESTAction{RESTActionStore, RESTActionUpdate}
	// RESTDelete 删除操作
	RESTDelete = []RESTAction{RESTActionDestory, RESTActionRestore}
	// RESTAdmin 管理操作
	RESTAdmin = []RESTAction{RESTActionList, RESTActionPage, RESTActionStore, RESTActionUpdate, RESTActionDestory, RESTActionRestore}
	// RESTAll 全部操作
	RESTAll = []RESTAction{RESTActionGetByID, RESTActionGet, RESTActionList, RESTActionPage, RESTActionStore, RESTActionUpdate, RESTActionDestory, RESTActionRestore}
	// ErrAction --
	ErrAction = errors.New("不支持Action")
)

// FilterFunc 过滤函数
type FilterFunc func(filters []Filter, r *http.Request) []Filter

// RESTHTTPService HTTP服务
type RESTHTTPService struct {
	repository *Repository
	filterFunc FilterFunc
}

// NewRESTHTTPService 新建HTTP服务
func NewRESTHTTPService(model interface{}) *RESTHTTPService {
	repository := NewRepository(model)
	return &RESTHTTPService{repository: repository}
}

// NewRESTHTTPServiceWithRepository 新建服务
func NewRESTHTTPServiceWithRepository(repository *Repository) *RESTHTTPService {
	return &RESTHTTPService{repository: repository}
}

// HandleFilterFunc 设置认证过滤参数函数
func (service *RESTHTTPService) HandleFilterFunc(handler FilterFunc) {
	service.filterFunc = handler
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
		handlerFunc = service.Get
		method = "GET"
	case RESTActionGet:
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

// Begin 开始事务
func (service *RESTHTTPService) Begin() *gorm.DB {
	return service.repository.Begin()
}

// Get 单个
func (service *RESTHTTPService) Get(w http.ResponseWriter, r *http.Request) {
	id, err := MuxParseVarID(r)
	if err != nil {
		NewFailResponse(err.Error()).WriteJSON(w)
		return
	}
	preloads := parseQueryPreloads(r)
	model := ReflectNew(service.repository.modelType)
	if err = service.repository.Get(id, model, preloads); err != nil {
		NewFailResponse(err.Error()).WriteJSON(w)
		return
	}
	NewSuccessResponse().
		AddData(reflectSingleKey(service.repository.modelValue), model).
		WriteJSON(w)
}

// First 单个
func (service *RESTHTTPService) First(w http.ResponseWriter, r *http.Request) {
	filterRequest, err := DecodeFilterRequest(r)
	if err != nil {
		NewFailResponse(err.Error()).WriteJSON(w)
		return
	}
	if nil != service.filterFunc {
		filterRequest.Filters = service.filterFunc(filterRequest.Filters, r)
	}
	model := ReflectNew(service.repository.modelType)
	if err = service.repository.First(filterRequest, model); err != nil {
		NewFailResponse(err.Error()).WriteJSON(w)
		return
	}
	NewSuccessResponse().
		AddData(reflectSingleKey(service.repository.modelValue), model).
		WriteJSON(w)
}

// List 列表
func (service *RESTHTTPService) List(w http.ResponseWriter, r *http.Request) {
	filterRequest, err := DecodeFilterRequest(r)
	if err != nil {
		NewFailResponse(err.Error()).WriteJSON(w)
		return
	}
	if nil != service.filterFunc {
		filterRequest.Filters = service.filterFunc(filterRequest.Filters, r)
	}
	list := ReflectNewSliceOfPtrTo(service.repository.modelType)
	if err = service.repository.List(filterRequest, list); err != nil {
		NewFailResponse(err.Error()).WriteJSON(w)
		return
	}
	NewSuccessResponse().
		AddData(reflectMultiKey(service.repository.modelValue), list).
		WriteJSON(w)
}

// Page 分页
func (service *RESTHTTPService) Page(w http.ResponseWriter, r *http.Request) {
	pageRequest, err := DecodePageRequest(r)
	if err != nil {
		NewFailResponse(err.Error()).WriteJSON(w)
		return
	}
	if nil != service.filterFunc {
		pageRequest.Filters = service.filterFunc(pageRequest.Filters, r)
	}
	pageResponse := &PageResponse{}
	pageResponse.List = ReflectNewSliceOfPtrTo(service.repository.modelType)
	if err = service.repository.Page(pageRequest, pageResponse); err != nil {
		NewFailResponse(err.Error()).WriteJSON(w)
		return
	}
	NewSuccessResponse().
		AddData(reflectMultiKey(service.repository.modelValue), pageResponse.List).
		AddData("pagination", pageResponse.Pagination).
		WriteJSON(w)
}

// Store 保存
func (service *RESTHTTPService) Store(w http.ResponseWriter, r *http.Request) {
	model, err := DecodeModel(r, service.repository.modelType)
	if err != nil {
		NewFailResponse(err.Error()).WriteJSON(w)
		return
	}
	if err = service.repository.Store(model); err != nil {
		NewFailResponse(err.Error()).WriteJSON(w)
		return
	}
	NewSuccessResponse().
		AddData(reflectSingleKey(service.repository.modelValue), model).
		WriteJSON(w)
}

// Update 更新
func (service *RESTHTTPService) Update(w http.ResponseWriter, r *http.Request) {
	id, err := MuxParseVarID(r)
	if err != nil {
		NewFailResponse(err.Error()).WriteJSON(w)
		return
	}
	model, err := DecodeModel(r, service.repository.modelType)
	if err != nil {
		NewFailResponse(err.Error()).WriteJSON(w)
		return
	}
	if err = service.repository.Update(id, model); err != nil {
		NewFailResponse(err.Error()).WriteJSON(w)
		return
	}
	NewSuccessResponse().
		AddData(reflectSingleKey(service.repository.modelValue), model).
		WriteJSON(w)
}

// Destory 删除
func (service *RESTHTTPService) Destory(w http.ResponseWriter, r *http.Request) {
	id, err := MuxParseVarID(r)
	if err != nil {
		NewFailResponse(err.Error()).WriteJSON(w)
		return
	}
	model := ReflectNew(service.repository.modelType)
	if err = service.repository.Destroy(id, model); err != nil {
		NewFailResponse(err.Error()).WriteJSON(w)
		return
	}
	NewSuccessResponse().
		AddData(reflectSingleKey(service.repository.modelValue), model).
		WriteJSON(w)
}

// Restore 恢复
func (service *RESTHTTPService) Restore(w http.ResponseWriter, r *http.Request) {
	id, err := MuxParseVarID(r)
	if err != nil {
		NewFailResponse(err.Error()).WriteJSON(w)
		return
	}
	model := ReflectNew(service.repository.modelType)
	if err = service.repository.Restore(id, model); err != nil {
		NewFailResponse(err.Error()).WriteJSON(w)
		return
	}
	NewSuccessResponse().
		AddData(reflectSingleKey(service.repository.modelValue), model).
		WriteJSON(w)
}
