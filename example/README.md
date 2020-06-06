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
+ HTTP
  + [http] [GET             ] /api/example/{id:[0-9]+}                                     PanicResponse, ReadMiddleware, TimeLogger
  + [http] [POST            ] /api/example/first                                           PanicResponse, ReadMiddleware, TimeLogger
  + [http] [POST            ] /api/example/list                                            PanicResponse, ReadMiddleware, TimeLogger
  + [http] [POST            ] /api/example/page                                            PanicResponse, ReadMiddleware, TimeLogger
  + [http] [POST            ] /api/example                                                 PanicResponse, WriteMiddleware, TimeLogger
  + [http] [PUT, POST       ] /api/example/{id:[0-9]+}                                     PanicResponse, WriteMiddleware, TimeLogger
  + [http] [PATCH, PUT, POST] /api/example/{id:[0-9]+}/fields                              PanicResponse, WriteMiddleware, TimeLogger
  + [http] [DELETE          ] /api/example/{id:[0-9]+}/remove                              PanicResponse, DeleteMiddleware, TimeLogger
  + [http] [DELETE          ] /api/example/{id:[0-9]+}/restore                             PanicResponse, DeleteMiddleware, TimeLogger
  + [http] [DELETE          ] /api/example/{id:[0-9]+}/destroy                             PanicResponse, DeleteMiddleware, TimeLogger
  + [http] [GET             ] /api/example_action                                          PanicResponse, ExampleMiddleware, TimeLogger
  + [http] [POST            ] /api/example_action                                          PanicResponse, ExampleMiddleware, TimeLogger
+ RPC
  + [  ws] /api/ws/once
  + [  ws] /api/ws/echo
+ WS
  + [ rpc] ExampleRPCService Get(string, *example.Example); List(gglmm.FilterRequest, *[]example.Example)