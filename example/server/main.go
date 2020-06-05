package main

import (
	"log"
	"net/http"

	"github.com/weihongguo/gglmm"
	redis "github.com/weihongguo/gglmm-redis"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	gglmm.RegisterGormDB("mysql", "example:123456@(127.0.0.1:3306)/example?charset=utf8mb4&parseTime=true&loc=UTC", 10, 5, 600)
	defer gglmm.CloseGormDB()

	redisCacher := redis.NewCacher("tcp", "127.0.0.1:6379", 5, 10, 3, 30)
	defer redisCacher.Close()
	gglmm.RegisterCacher(redisCacher)

	gglmm.BasePath("/api")

	// 认证中间间
	// authenticationMiddlerware := gglmm.authenticationMiddlerware("example")

	gglmm.UseTimeLogger(true)

	exampleService := NewExampleService()
	exampleService.
		HandleBeforeCreateFunc(beforeCreate).
		HandleBeforeUpdateFunc(beforeUpdate)

	readMiddleware := gglmm.Middleware{
		Name: "ReadMiddleware",
		Func: middlewareFunc,
	}
	gglmm.HandleHTTP("/example", exampleService).
		Action(readMiddleware, gglmm.ReadActions)

	writeMiddleware := gglmm.Middleware{
		Name: "WriteMiddleware",
		Func: middlewareFunc,
	}
	gglmm.HandleHTTP("/example", exampleService).
		Action(writeMiddleware, gglmm.WriteActions)

	deleteMiddleware := gglmm.Middleware{
		Name: "DeleteMiddleware",
		Func: middlewareFunc,
	}
	gglmm.HandleHTTP("/example", exampleService).
		Action(deleteMiddleware, gglmm.DeleteActions)

	exampleMiddleware := gglmm.Middleware{
		Name: "ExampleMiddleware",
		Func: middlewareFunc,
	}
	gglmm.HandleHTTPAction("/example_action", exampleService.ExampleAction, "GET").
		Middleware(exampleMiddleware)
	gglmm.HandleHTTPAction("/example_action", ExampleAction, "POST").
		Middleware(exampleMiddleware)

	gglmm.HandleWS("/ws/once", OnceWSHandler)
	gglmm.HandleWS("/ws/echo", EchoWSHandler)

	gglmm.RegisterRPC(NewExampleRPCService())

	gglmm.ListenAndServe(":10000")
}

func middlewareFunc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("middleware before")
		next.ServeHTTP(w, r)
		log.Printf("middleware after")
		return
	})
}

func beforeCreate(model interface{}) (interface{}, error) {
	log.Printf("%#v\n", model)
	example, ok := model.(*Example)
	if !ok {
		return nil, gglmm.ErrModelType
	}
	example.StringValue = "string"
	return example, nil
}

func beforeUpdate(model interface{}, id int64) (interface{}, int64, error) {
	example, ok := model.(*Example)
	if !ok {
		return nil, 0, gglmm.ErrModelType
	}
	example.StringValue = "string"
	return example, id, nil
}
