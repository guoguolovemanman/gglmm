package gglmm

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

// FilterFieldDeleted --
const FilterFieldDeleted = "deleted"

var (
	// FilterValueAll --
	FilterValueAll = ConfigString{Value: "all", Name: "所有"}
	// FilterValueDeleted --
	FilterValueDeleted = ConfigString{Value: "deleted", Name: "已删除"}
)

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
	ID       int64    `json:"id"`
	Preloads []string `json:"preloads"`
}

// FilterRequest 分页请求
type FilterRequest struct {
	Filters  []Filter `json:"filters"`
	Order    string   `json:"order"`
	Preloads []string `json:"preloads"`
}

// AddFilter 添加过滤条件
func (request *FilterRequest) AddFilter(field string, operate string, value interface{}) {
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
	*FilterRequest
	Pagination
}
