package gglmm

import "errors"

// Err --
var (
	ErrRequest = errors.New("请求参数错误")
)

// Filter
const (
	FilterOperateEqual        = "="
	FilterOperateNotEqual     = "<>"
	FilterOperateGreaterThan  = ">"
	FilterOperateGreaterEqual = ">="
	FilterOperateLessThan     = "<"
	FilterOperateLessEqual    = "<="
	FilterOperateLike         = "like"
	FilterOperateIn           = "in"
	FilterOperateBetween      = "between"
	FilterSeparator           = "|"
)

// Filter
var (
	FilterFieldDeleted = "deleted"
	FilterValueAll     = ConfigString{Value: "all", Name: "所有"}
	FilterValueDeleted = ConfigString{Value: "deleted", Name: "已删除"}
)

// Page
const (
	DefaultPageSize = 15
	FirstPageIndex  = 1
)

// Filter 过滤参数
type Filter struct {
	Field   string      `json:"field"`
	Operate string      `json:"operate"`
	Value   interface{} `json:"value"`
}

// NewFilter --
func NewFilter(field string, operate string, value interface{}) *Filter {
	return &Filter{
		Field:   field,
		Operate: operate,
		Value:   value,
	}
}

// Check --
func (filter Filter) Check() bool {
	if filter.Field == "" {
		return false
	}
	if filter.Operate == "" {
		return false
	}
	if filter.Value == nil {
		return false
	}
	return true
}

// IDRequest --
type IDRequest struct {
	ID       uint64   `json:"id"`
	Preloads []string `json:"preloads"`
}

// FilterRequest 分页请求
type FilterRequest struct {
	Filters  []*Filter `json:"filters"`
	Preloads []string  `json:"preloads"`
	Order    string    `json:"order"`
}

// AddFilter 添加过滤条件
func (request *FilterRequest) AddFilter(field string, operate string, value interface{}) {
	if request.Filters == nil {
		request.Filters = make([]*Filter, 0)
	}
	request.Filters = append(request.Filters, NewFilter(field, operate, value))
}

// Pagination 分页
type Pagination struct {
	PageSize  int `json:"pageSize"`
	PageIndex int `json:"pageIndex"`
	Total     int `json:"total"`
}

// PageRequest 分页请求
type PageRequest struct {
	*FilterRequest
	Pagination
}
