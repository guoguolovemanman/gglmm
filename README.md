# gglmm
## `gglmm` 怎么理解？
+ `gg`：本人昵称首字母
+ `l`: love首字母
+ `mm`：爱人昵称首字母
## 依赖
+ github.com/gorilla/mux  路由
+ github.com/jinzhu/gorm  数据库
+ github.com/dgrijalva/jwt-go 认证
+ golang.org/x/crypto 密码
## 基本模型
```golang
type Model struct {
	ID        int64      `json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt"`
}
```
## 用法
+ **详见example**
+ 数据库 -- 依赖gorm库
```golang
func RegisterGormDBConfig(config ConfigDB)
func RegisterGormDB(dialect string, url string, maxOpen int, maxIdle int, connMaxLifetime time.Duration)
func CloseGormDB()
func DefaultGormDB() *GormDB
```
+ 缓存（gglmm-redis库提供了一个依赖redis的实现）
```golang
func RegisterCacher(cacherInstance Cacher)
func DefaultCacher() Cacher
```
+ HTTP
```golang
type HTTPHandler interface {
	Action(action Action) (*HTTPAction, error)
}

// 注册HTTPHandler（HTTPAction集合）
func HandleHTTP(path string, httpHandler HTTPHandler) *HTTPHandlerConfig

// 声明HTTPHandler的Action，本次声明的所有Middleware作用于所有Action
// param: Middleware | []Middleware | Action | []Action
// Middleware 按顺序作用
func (config *HTTPHandlerConfig) Action(params ...interface{}) *HTTPHandlerConfig

// 注册HTTPAction
func HandleHTTPAction(path string, handlerFunc http.HandlerFunc, methods ...string) *HTTPActionConfig

// 声明HTTPAction的Middleware
// Middleware 按顺序作用
func (config *HTTPActionConfig) Middleware(middlewares ...Middleware)
```
+ HTTPService 内部的`HTTPHandler`实现
```golang
func NewHTTPService(model interface{}) *HTTPService

// 模型自定义返回结果的Key，默认为[record, records]
func (model Model) ResponseKey() [2]string

// 模型自定义是否支持缓存，默认false
func (model Model) Cache() bool

// 根据Action名称注册HTTPAction
func (service *HTTPService) Action(action Action) (*HTTPAction, error)

// 内部实现了以下Action
const (
	// ActionGetByID 根据ID拉取单个
	ActionGetByID Action = "GetByID"
	// ActionFirst 根据条件拉取单个
	ActionFirst Action = "First"
	// ActionList 列表
	ActionList Action = "List"
	// ActionPage 分页
	ActionPage Action = "page"
	// ActionCreate 保存
	ActionCreate Action = "create"
	// ActionUpdate 更新整体
	ActionUpdate Action = "Update"
	// ActionUpdateFields 更新多个字段
	ActionUpdateFields Action = "UpdateFields"
	// ActionRemove 软删除
	ActionRemove Action = "Remove"
	// ActionRestore 恢复
	ActionRestore Action = "Resotre"
	// ActionDestory 硬删除
	ActionDestory Action = "Destory"
)

// GET basePath/resourcePath/{id:[0-9]+} 根据ID查询
func (service *HTTPService) GetByID(w http.ResponseWriter, r *http.Request)

// POST basePaht/resourcePaht/fist 根据条件查询第一个
func (service *HTTPService) First(w http.ResponseWriter, r *http.Request)

// POST basePaht/resourcePaht/list 根据条件查询，输出列表
func (service *HTTPService) List(w http.ResponseWriter, r *http.Request)

// POST basePaht/resourcePaht/page 根据条件查询，输出分页
func (service *HTTPService) Page(w http.ResponseWriter, r *http.Request)

// POST basePaht/resourcePaht 保存
func (service *HTTPService) Store(w http.ResponseWriter, r *http.Request)

// PUT basePaht/resourcePaht/{id:[0-9]+} 更新整体
func (service *HTTPService) Update(w http.ResponseWriter, r *http.Request)

// PATCH basePaht/resourcePaht/{id:[0-9]+} 更新部分字段
func (service *HTTPService) UpdateFields(w http.ResponseWriter, r *http.Request)

// DELETE basePaht/resourcePaht/{id:[0-9]+}/remove 软删除
func (service *HTTPService) Remove(w http.ResponseWriter, r *http.Request)

// DELETE basePaht/resourcePaht/{id:[0-9]+}/restore 恢复软删除
func (service *HTTPService) Restore(w http.ResponseWriter, r *http.Request)

// DELETE basePaht/resourcePaht/{id:[0-9]+}/destroy 硬删除
func (service *HTTPService) Destory(w http.ResponseWriter, r *http.Request)
```
+ RPC
```golang
type RPCHandler interface {
	Actions(cmd string, actions *[]*RPCAction) error
}

// 注册RPCHandler
func RegisterRPC(rpcHandler RPCHandler) *RPCHandlerConfig

// 注册RPCHandler，指定名称
func RegisterRPCName(name string, rpcHandler RPCHandler) *RPCHandlerConfig
```