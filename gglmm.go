package gglmm

import (
	"log"
	"net/http"
	"net/rpc"

	"github.com/gorilla/mux"
)

var basePath string = ""

var usePanicResponse = true
var middlewarePanicResponse = MiddlewarePanicResponse()

var useTimeLogger = false
var middlewareTimeLogger = MiddlewareTimeLogger()

// BasePath 基础路径
func BasePath(path string) {
	basePath = path
}

// UsePanicResponse --
func UsePanicResponse(use bool) {
	usePanicResponse = use
}

// UseTimeLogger --
func UseTimeLogger(use bool) {
	useTimeLogger = use
}

// ListenAndServe 监听并服务
func ListenAndServe(address string) {
	log.Println("listen on: " + address)

	router := mux.NewRouter()
	handleHTTP(router)
	handleHTTPAction(router)
	http.Handle("/", router)

	handleWS()

	registerRPC()
	rpc.HandleHTTP()

	err := http.ListenAndServe(address, nil)
	if err != nil {
		panic(err)
	}
}

// ListenAndServeConfig 监听并服务
func ListenAndServeConfig(config ConfigAPI) {
	if !config.Check() {
		log.Fatal("ConfigAPI invalid")
	}
	ListenAndServe(config.Address)
}
