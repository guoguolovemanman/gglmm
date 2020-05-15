# gglmm
## `gglmm` 再怎么理解？
+ `gg`：本人昵称首字母
+ `l`: love首字母
+ `mm`：本人爱人昵称首字母
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
## HTTP接口
```golang
type HTTPHandler interface {
	Action(action string) (*HTTPAction, error)
}
```
## RPC接口
```golang
type RPCHandler interface {
	Actions(cmd string, actions *[]*RPCAction) error
}
```
## 用法
+ **详见example**
+ 数据库 -- Gorm
```golang
func RegisterGormDBConfig(config ConfigDB)
func RegisterGormDB(dialect string, url string, maxOpen int, maxIdle int, connMaxLifetime time.Duration)
func CloseGormDB()
func DefaultGormDB() *GormDB
```
+ 缓存（参考gglmm-redis库）
```golang
func RegisterCacher(cacherInstance Cacher)
func DefaultCacher() Cacher
```
+ HTTP
```golang
func HandleHTTP(path string, httpHandler HTTPHandler) *HTTPHandlerConfig
func HandleHTTPAction(path string, handlerFunc http.HandlerFunc, methods ...string) *HTTPActionConfig
```
+ HTTPService 实现了 HTTPHandler 接口
```golang
func NewHTTPService(model interface{}) *HTTPService
func (service *HTTPService) Action(action string) (*HTTPAction, error)

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
func RegisterRPC(rpcHandler RPCHandler) *RPCHandlerConfig
func RegisterRPCName(name string, rpcHandler RPCHandler) *RPCHandlerConfig
```