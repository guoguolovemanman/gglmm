package main

import (
	"log"
	"net/http"

	"github.com/weihongguo/gglmm"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// Example --
type Example struct {
	gglmm.Model
	Example string `json:"example"`
}

// RPCExampleService --
type RPCExampleService struct {
}

// Actions --
func (service *RPCExampleService) Actions(cmd string, actions *[]string) error {
	*actions = []string{"RPCAction1(string, *string)", "RPCAction1(string, *string)"}
	return nil
}

// RPCAction1 --
func (service *RPCExampleService) RPCAction1(req string, res *string) error {
	*res = "RPCAction1"
	return nil
}

// RPCAction2 --
func (service *RPCExampleService) RPCAction2(req string, res *string) error {
	*res = "RPCAction2"
	return nil
}

func main() {
	gglmm.RegisterGormDB("mysql", "example:123456@(127.0.0.1:3306)/example?charset=utf8mb4&parseTime=true&loc=UTC", 10, 5, 600)
	defer gglmm.CloseGormDB()

	gglmm.RegisterBasePath("/api/example")

	// 登录态中间件请参考gglmm-account
	// 缓存可以参考gglmm-redis里的Cacher实现

	gglmm.RegisterHTTPHandler(gglmm.NewRESTHTTPService(Example{}), "/example").
		Middleware(gglmm.Middleware{
			Name: "example",
			Func: middlewareFunc,
		}).
		RESTAction(gglmm.RESTAll)

	gglmm.RegisterRPCHandler(&RPCExampleService{}, "RPCExampleService")

	gglmm.ListenAndServe(":10000")
}

func middlewareFunc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("before")
		next.ServeHTTP(w, r)
		log.Printf("after")
		return
	})
}
