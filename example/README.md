# gglmm-example
## 模型
```golang
type Example struct {
	gglmm.Model
	IntValue    int     `json:"intValue"`
	FloatValue  float64 `json:"floatValue"`
	StringValue string  `json:"stringValue"`
}
```
## Server
+ 用法
```shell
go build
./gglmm-example-server
```
+ HTTP
```golang
// 例子Service，基于框架的HTTPService
type ExampleService struct {
	*gglmm.HTTPService
}
exampleService := NewExampleService()

// 读
readMiddleware := gglmm.Middleware{
  Name: "ReadMiddleware",
  Func: middlewareFunc,
}
gglmm.HandleHTTP("/example", exampleService).Action(readMiddleware, gglmm.ReadActions)

// 写
writeMiddleware := gglmm.Middleware{
  Name: "WriteMiddleware",
  Func: middlewareFunc,
}
gglmm.HandleHTTP("/example", exampleService).Action(writeMiddleware, gglmm.WriteActions)

// 删
deleteMiddleware := gglmm.Middleware{
  Name: "DeleteMiddleware",
  Func: middlewareFunc,
}
gglmm.HandleHTTP("/example", exampleService).Action(deleteMiddleware, gglmm.DeleteActions)

// 自定以
exampleMiddleware := gglmm.Middleware{
  Name: "ExampleMiddleware",
  Func: middlewareFunc,
}
// Service内部函数
gglmm.HandleHTTPAction("/example_action", exampleService.ExampleAction, "GET").Middleware(exampleMiddleware)
// 直接函数
gglmm.HandleHTTPAction("/example_action", ExampleAction, "POST").Middleware(exampleMiddleware)
```
+ RPC
```golang
type ExampleRPCService struct {
	gormDB *gglmm.GormDB
}
func (service *ExampleRPCService) Actions(cmd string, actions *[]*gglmm.RPCAction) error {
	*actions = append(*actions, []*gglmm.RPCAction{
		gglmm.NewRPCAction("Get", "string", "*example.Example"),
		gglmm.NewRPCAction("List", "gglmm.FilterRequest", "*[]example.Example"),
	}...)
	return nil
}
```
+ WS
```golang
// 应搭一次
gglmm.HandleWS("/ws/once", OnceWSHandler)
// 回声
gglmm.HandleWS("/ws/echo", EchoWSHandler)
```