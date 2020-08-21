package main

import (
	example "gglmm-example"
	"log"
	"net/http"

	"github.com/weihongguo/gglmm"
	auth "github.com/weihongguo/gglmm-auth"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	gglmm.RegisterGormDB("mysql", "gglmm_example:123456@(127.0.0.1:3306)/gglmm_example?charset=utf8mb4&parseTime=true&loc=UTC", 10, 5, 600)
	defer gglmm.CloseGormDB()

	gglmm.BasePath("/api")

	// 认证中间间
	// authenticationMiddlerware := gglmm.authenticationMiddlerware("example")

	jwtAuthExample := auth.MiddlewareJWTAuthChecker("example")
	permission := gglmm.MiddlewarePermissionChecker(checkPermission)

	gglmm.UseTimeLogger(true, 300)

	gglmm.HandleHTTPAction("/login", LoginAction(3600, "example"), "GET", "POST")

	exampleService := NewExampleService()
	exampleService.
		HandleBeforeCreateFunc(beforeCreate).
		HandleBeforeUpdateFunc(beforeUpdate)

	readMiddleware := gglmm.Middleware{
		Name: "ReadMiddleware",
		Func: middlewareFunc,
	}
	gglmm.HandleHTTP("/example", exampleService).
		Action(jwtAuthExample, permission, readMiddleware, gglmm.ReadActions)

	writeMiddleware := gglmm.Middleware{
		Name: "WriteMiddleware",
		Func: middlewareFunc,
	}
	gglmm.HandleHTTP("/example", exampleService).
		Action(jwtAuthExample, permission, writeMiddleware, gglmm.WriteActions)

	deleteMiddleware := gglmm.Middleware{
		Name: "DeleteMiddleware",
		Func: middlewareFunc,
	}
	gglmm.HandleHTTP("/example", exampleService).
		Action(jwtAuthExample, permission, deleteMiddleware, gglmm.DeleteActions)

	exampleMiddleware := &gglmm.Middleware{
		Name: "ExampleMiddleware",
		Func: middlewareFunc,
	}
	gglmm.HandleHTTPAction("/example_action", exampleService.ExampleAction, "GET").
		Middleware(jwtAuthExample, permission, exampleMiddleware)
	gglmm.HandleHTTPAction("/example_action", ExampleAction, "POST").
		Middleware(jwtAuthExample, permission, exampleMiddleware)

	gglmm.HandleWS("/ws/once", OnceWSHandler)
	gglmm.HandleWS("/ws/echo", EchoWSHandler)

	gglmm.RegisterRPC(NewExampleRPCService())

	gglmm.ListenAndServe(":10000")
}

func checkPermission(r *http.Request) error {
	authSubject, err := auth.SubjectFrom(r)
	if err != nil {
		return err
	}
	log.Println(r.Method, r.URL.Path, authSubject.UserType, authSubject.UserID)
	return nil
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
	example, ok := model.(*example.Example)
	if !ok {
		return nil, gglmm.ErrModelType
	}
	example.StringValue = "string"
	return example, nil
}

func beforeUpdate(model interface{}) (interface{}, error) {
	example, ok := model.(*example.Example)
	if !ok {
		return nil, gglmm.ErrModelType
	}
	example.StringValue = "string"
	return example, nil
}
