package gglmm

import (
	"errors"
	"net/http"
	"reflect"
)

// Err
var (
	ErrAction            = errors.New("不支持Action")
	ErrModelType         = errors.New("模型类型错误")
	ErrModelCanNotDelete = errors.New("模型不可删除")
	ErrModelCanNotUpdate = errors.New("模型不可更新")
)

// Action --
type Action string

// Action --
const (
	ActionGetByID Action = "GetByID"
	ActionFirst   Action = "First"
	ActionAdmin   Action = "Admin"
	ActionList    Action = "List"
	ActionPage    Action = "Page"
	ActionCreate  Action = "Create"
	ActionStore   Action = "Store"
	ActionEdit    Action = "Edit"
	ActionUpdate  Action = "Update"
	ActionRemove  Action = "Remove"
	ActionRestore Action = "Resotre"
	ActionDestory Action = "Destory"
)

// IDRegexp ID正则表达式
const IDRegexp = "{id:[0-9]+}"

var (
	// ReadActions 读Action
	ReadActions = []Action{ActionGetByID, ActionFirst, ActionList, ActionPage}
	// WriteActions 写Action
	WriteActions = []Action{ActionStore, ActionUpdate}
	// DeleteActions 删除Action
	DeleteActions = []Action{ActionRemove, ActionRestore, ActionDestory}
)

// FilterFunc 过滤函数
type FilterFunc func([]*Filter, *http.Request) []*Filter

// BeforeCreateFunc 保存前调用
type BeforeCreateFunc func(interface{}, *http.Request) (interface{}, error)

// BeforeUpdateFunc 更新前调用
type BeforeUpdateFunc func(interface{}, *http.Request) (interface{}, error)

// BeforeDeleteFunc 删除前调用
type BeforeDeleteFunc func(interface{}, *http.Request) (interface{}, error)

// HTTPService HTTP服务
type HTTPService struct {
	gglmmDB    *DB
	entityType reflect.Type
	keys       [2]string

	filterFunc       FilterFunc
	beforeCreateFunc BeforeCreateFunc
	beforeUpdateFunc BeforeUpdateFunc
	beforeDeleteFunc BeforeDeleteFunc
}

// NewHTTPService 新建HTTP服务
func NewHTTPService(entity Entity, keys [2]string) *HTTPService {
	return &HTTPService{
		gglmmDB:    NewDB(),
		entityType: reflect.TypeOf(entity),
		keys:       keys,
	}
}

// HandleFilterFunc 设置过滤参数函数
func (service *HTTPService) HandleFilterFunc(handler FilterFunc) *HTTPService {
	service.filterFunc = handler
	return service
}

// HandleBeforeCreateFunc 设置保存前执行函数
func (service *HTTPService) HandleBeforeCreateFunc(handler BeforeCreateFunc) *HTTPService {
	service.beforeCreateFunc = handler
	return service
}

// HandleBeforeUpdateFunc 设置更新前执行函数
func (service *HTTPService) HandleBeforeUpdateFunc(handler BeforeUpdateFunc) *HTTPService {
	service.beforeUpdateFunc = handler
	return service
}

// HandleBeforeDeleteFunc 设置更新前执行函数
func (service *HTTPService) HandleBeforeDeleteFunc(handler BeforeDeleteFunc) *HTTPService {
	service.beforeDeleteFunc = handler
	return service
}

// Action --
func (service *HTTPService) Action(action Action) (*HTTPAction, error) {
	var path string
	var handlerFunc http.HandlerFunc
	var methods []string
	switch action {
	case ActionGetByID:
		path = "/" + IDRegexp
		handlerFunc = service.GetByID
		methods = []string{"GET"}
	case ActionFirst:
		path = "/first"
		handlerFunc = service.First
		methods = []string{"POST"}
	case ActionList:
		path = "/list"
		handlerFunc = service.List
		methods = []string{"POST"}
	case ActionPage:
		path = "/page"
		handlerFunc = service.Page
		methods = []string{"POST"}
	case ActionStore:
		handlerFunc = service.Store
		methods = []string{"POST"}
	case ActionUpdate:
		path = "/" + IDRegexp
		handlerFunc = service.Update
		methods = []string{"PUT", "POST"}
	case ActionRemove:
		path = "/" + IDRegexp + "/remove"
		handlerFunc = service.Remove
		methods = []string{"DELETE"}
	case ActionRestore:
		path = "/" + IDRegexp + "/restore"
		handlerFunc = service.Restore
		methods = []string{"DELETE"}
	case ActionDestory:
		path = "/" + IDRegexp + "/destroy"
		handlerFunc = service.Destory
		methods = []string{"DELETE"}
	}
	if handlerFunc != nil {
		return NewHTTPAction(path, handlerFunc, methods...), nil
	}
	return nil, ErrAction
}

// GetByID 单个
func (service *HTTPService) GetByID(w http.ResponseWriter, r *http.Request) {
	idRequest := IDRequest{}
	if err := DecodeBody(r, &idRequest); err != nil {
		FailResponse(NewErrFileLine(err)).JSON(w)
		return
	}
	entity := reflect.New(service.entityType).Interface()
	if err := service.gglmmDB.First(entity, idRequest); err != nil {
		FailResponse(NewErrFileLine(err)).JSON(w)
		return
	}
	OkResponse().
		AddData(service.keys[0], entity).
		JSON(w)
}

// First 单个
func (service *HTTPService) First(w http.ResponseWriter, r *http.Request) {
	filterRequest := FilterRequest{}
	if err := DecodeBody(r, &filterRequest); err != nil {
		FailResponse(NewErrFileLine(err)).JSON(w)
		return
	}
	if service.filterFunc != nil {
		filterRequest.Filters = service.filterFunc(filterRequest.Filters, r)
	}
	entity := reflect.New(service.entityType).Interface()
	if err := service.gglmmDB.First(entity, filterRequest); err != nil {
		FailResponse(NewErrFileLine(err)).JSON(w)
		return
	}
	OkResponse().
		AddData(service.keys[0], entity).
		JSON(w)
}

// List 列表
func (service *HTTPService) List(w http.ResponseWriter, r *http.Request) {
	filterRequest := FilterRequest{}
	if err := DecodeBody(r, &filterRequest); err != nil {
		FailResponse(NewErrFileLine(err)).JSON(w)
		return
	}
	if service.filterFunc != nil {
		filterRequest.Filters = service.filterFunc(filterRequest.Filters, r)
	}
	entities := reflect.New(reflect.SliceOf(service.entityType)).Interface()
	if err := service.gglmmDB.List(entities, &filterRequest); err != nil {
		FailResponse(NewErrFileLine(err)).JSON(w)
		return
	}
	OkResponse().
		AddData(service.keys[1], entities).
		JSON(w)
}

// Page 分页
func (service *HTTPService) Page(w http.ResponseWriter, r *http.Request) {
	pageRequest := PageRequest{}
	if err := DecodeBody(r, &pageRequest); err != nil {
		FailResponse(NewErrFileLine(err)).JSON(w)
		return
	}
	if service.filterFunc != nil {
		pageRequest.Filters = service.filterFunc(pageRequest.Filters, r)
	}
	pageResponse := &PageResponse{}
	pageResponse.List = reflect.New(reflect.SliceOf(service.entityType)).Interface()
	if err := service.gglmmDB.Page(pageResponse, &pageRequest); err != nil {
		FailResponse(NewErrFileLine(err)).JSON(w)
		return
	}
	OkResponse().
		AddData(service.keys[1], pageResponse.List).
		AddData("pagination", pageResponse.Pagination).
		JSON(w)
}

// Store 保存
func (service *HTTPService) Store(w http.ResponseWriter, r *http.Request) {
	entity := reflect.New(service.entityType).Interface()
	err := DecodeBody(r, entity)
	if err != nil {
		FailResponse(NewErrFileLine(err)).JSON(w)
		return
	}
	if service.beforeCreateFunc != nil {
		entity, err = service.beforeCreateFunc(entity, r)
		if err != nil {
			FailResponse(NewErrFileLine(err)).JSON(w)
			return
		}
	}
	if err := service.gglmmDB.Create(entity); err != nil {
		FailResponse(NewErrFileLine(err)).JSON(w)
		return
	}
	OkResponse().
		AddData(service.keys[0], entity).
		JSON(w)
}

// Update 更新整体
func (service *HTTPService) Update(w http.ResponseWriter, r *http.Request) {
	id, err := PathVarID(r)
	if err != nil {
		FailResponse(NewErrFileLine(err)).JSON(w)
		return
	}
	entity := reflect.New(service.entityType).Interface()
	if err = DecodeBody(r, entity); err != nil {
		FailResponse(NewErrFileLine(err)).JSON(w)
		return
	}
	SetPrimaryKeyValue(entity, id)
	if service.beforeUpdateFunc != nil {
		entity, err = service.beforeUpdateFunc(entity, r)
		if err != nil {
			FailResponse(NewErrFileLine(err)).JSON(w)
			return
		}
	}
	if err = service.gglmmDB.Update(entity); err != nil {
		FailResponse(NewErrFileLine(err)).JSON(w)
		return
	}
	OkResponse().
		AddData(service.keys[0], entity).
		JSON(w)
}

// Remove 软删除
func (service *HTTPService) Remove(w http.ResponseWriter, r *http.Request) {
	id, err := PathVarID(r)
	if err != nil {
		FailResponse(NewErrFileLine(err)).JSON(w)
		return
	}
	entity := reflect.New(service.entityType).Interface()
	if service.beforeDeleteFunc != nil {
		if err := service.gglmmDB.First(entity, id); err != nil {
			FailResponse(NewErrFileLine(err)).JSON(w)
			return
		}
		if _, err := service.beforeDeleteFunc(entity, r); err != nil {
			FailResponse(NewErrFileLine(err)).JSON(w)
			return
		}
	} else {
		SetPrimaryKeyValue(entity, id)
	}
	if err = service.gglmmDB.Remove(entity); err != nil {
		FailResponse(NewErrFileLine(err)).JSON(w)
		return
	}
	OkResponse().
		AddData(service.keys[0], entity).
		JSON(w)
}

// Restore 恢复
func (service *HTTPService) Restore(w http.ResponseWriter, r *http.Request) {
	id, err := PathVarID(r)
	if err != nil {
		FailResponse(NewErrFileLine(err)).JSON(w)
		return
	}
	entity := reflect.New(service.entityType).Interface()
	SetPrimaryKeyValue(entity, id)
	if err = service.gglmmDB.Restore(entity); err != nil {
		FailResponse(NewErrFileLine(err)).JSON(w)
		return
	}
	OkResponse().
		AddData(service.keys[0], entity).
		JSON(w)
}

// Destory 直接删除
func (service *HTTPService) Destory(w http.ResponseWriter, r *http.Request) {
	id, err := PathVarID(r)
	if err != nil {
		FailResponse(NewErrFileLine(err)).JSON(w)
		return
	}
	entity := reflect.New(service.entityType).Interface()
	if service.beforeDeleteFunc != nil {
		if err := service.gglmmDB.First(entity, id); err != nil {
			FailResponse(NewErrFileLine(err)).JSON(w)
			return
		}
		if _, err := service.beforeDeleteFunc(entity, r); err != nil {
			FailResponse(NewErrFileLine(err)).JSON(w)
			return
		}
	} else {
		SetPrimaryKeyValue(entity, id)
	}
	if err = service.gglmmDB.Destroy(entity); err != nil {
		FailResponse(NewErrFileLine(err)).JSON(w)
		return
	}
	OkResponse().JSON(w)
}
