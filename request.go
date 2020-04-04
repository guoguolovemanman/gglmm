package gglmm

import (
	"errors"
)

// FilterOperate 过滤操作
type FilterOperate string

const (
	// FilterOperateEqual 等于
	FilterOperateEqual FilterOperate = "="
	// FilterOperateNotEqual 等于
	FilterOperateNotEqual FilterOperate = "<>"
	// FilterOperateGreaterThan 大于
	FilterOperateGreaterThan FilterOperate = ">"
	// FilterOperateGreaterEqual 大于等于
	FilterOperateGreaterEqual FilterOperate = ">="
	// FilterOperateLessThan 小于
	FilterOperateLessThan FilterOperate = "<"
	// FilterOperateLessEqual 小于等于
	FilterOperateLessEqual FilterOperate = "<="
	// FilterOperateLike like模糊匹配
	FilterOperateLike FilterOperate = "like"
	// FilterOperateIn in查询
	FilterOperateIn FilterOperate = "in"
	// FilterOperateBetween between查询
	FilterOperateBetween FilterOperate = "between"
	// FilterSeparator 参数风隔符
	FilterSeparator = "|"
	// DefaultPageSize 默认每页大小
	DefaultPageSize = 15
	// FirstPageIndex 第一页
	FirstPageIndex = 1
)

// ErrRequestParam 请求参数错误
var ErrRequestParam = errors.New("请求参数错误！")

// IDRequest --
type IDRequest struct {
	ID       int64    `json:"id"`
	Preloads []string `json:"preloads"`
}

// Filter 过滤参数
type Filter struct {
	Field   string        `json:"field"`
	Operate FilterOperate `json:"operate"`
	Value   interface{}   `json:"value"`
}

// FilterRequest 分页请求
type FilterRequest struct {
	Filters  []Filter `json:"filters"`
	Order    string   `json:"order"`
	Preloads []string `json:"preloads"`
}

// AddFilter 添加过滤条件
func (request *FilterRequest) AddFilter(field string, operate FilterOperate, value interface{}) {
	if request.Filters == nil {
		request.Filters = make([]Filter, 0)
	}
	filter := Filter{
		Field:   field,
		Operate: operate,
		Value:   value,
	}
	request.Filters = append(request.Filters, filter)
}

// Pagination 分页
type Pagination struct {
	PageSize  int `json:"pageSize"`
	PageIndex int `json:"pageIndex"`
	Total     int `json:"total"`
}

// PageRequest 分页请求
type PageRequest struct {
	FilterRequest
	Pagination
}
