package gglmm

import (
	"fmt"
	"log"
	"net/http"
	"net/rpc"
	"reflect"
	"strings"

	"github.com/gorilla/mux"
)

var basePath string = ""
var usePanicResponse = true
var panicResponseMiddleware = PanicResponse()
var useTimeLogger = false
var timeLoggerMiddleware = TimeLogger()
var httpHandlerConfigs []*HTTPHandlerConfig = nil
var httpActionConfigs []*HTTPActionConfig = nil
var rpcHandlerConfigs []*RPCHandlerConfig = nil

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

// HandleHTTP 注册HTTPHandler
// path 路径
// httpHandler 处理者
func HandleHTTP(path string, httpHandler HTTPHandler) *HTTPHandlerConfig {
	if httpHandlerConfigs == nil {
		httpHandlerConfigs = make([]*HTTPHandlerConfig, 0)
	}
	config := &HTTPHandlerConfig{
		path:        path,
		httpHandler: httpHandler,
	}
	httpHandlerConfigs = append(httpHandlerConfigs, config)
	return config
}

// HandleHTTPAction 注册HandlerFunc
// path 路径
// methods 方法
func HandleHTTPAction(path string, handlerFunc http.HandlerFunc, methods ...string) *HTTPActionConfig {
	if httpActionConfigs == nil {
		httpActionConfigs = make([]*HTTPActionConfig, 0)
	}
	if methods == nil {
		methods = []string{"GET"}
	}
	config := &HTTPActionConfig{
		httpAction: HTTPAction{
			path:        path,
			handlerFunc: handlerFunc,
			methods:     methods,
		},
	}
	httpActionConfigs = append(httpActionConfigs, config)
	return config
}

// RegisterRPC 注册RPCHandler
// rpcHandler 处理者
func RegisterRPC(rpcHandler RPCHandler) *RPCHandlerConfig {
	return RegisterRPCName("", rpcHandler)
}

// RegisterRPCName 注册RPCHandler
// name 名称
// rpcHandler 处理者
func RegisterRPCName(name string, rpcHandler RPCHandler) *RPCHandlerConfig {
	if rpcHandlerConfigs == nil {
		rpcHandlerConfigs = make([]*RPCHandlerConfig, 0)
	}
	config := &RPCHandlerConfig{
		Name:       name,
		RPCHandler: rpcHandler,
	}
	rpcHandlerConfigs = append(rpcHandlerConfigs, config)
	return config
}

func handleHTTP(router *mux.Router) {
	if httpHandlerConfigs == nil || len(httpHandlerConfigs) == 0 {
		return
	}
	for _, config := range httpHandlerConfigs {
		subrouter := router.PathPrefix(basePath).Subrouter()
		for _, middlewareAcion := range config.middlewareActions {
			middlewares := make([]string, 0)
			if usePanicResponse {
				subrouter.Use(mux.MiddlewareFunc(panicResponseMiddleware.Func))
				middlewares = append(middlewares, panicResponseMiddleware.Name)
			}
			for _, middleware := range middlewareAcion.middlewares {
				subrouter.Use(mux.MiddlewareFunc(middleware.Func))
				middlewares = append(middlewares, middleware.Name)
			}
			if useTimeLogger {
				subrouter.Use(mux.MiddlewareFunc(timeLoggerMiddleware.Func))
				middlewares = append(middlewares, timeLoggerMiddleware.Name)
			}
			for _, action := range middlewareAcion.actions {
				httpAction, err := config.httpHandler.Action(action)
				if err != nil {
					log.Println(err)
				} else if httpAction.handlerFunc != nil {
					path := config.path + httpAction.path
					handleHTTPFunc(subrouter, path, httpAction.handlerFunc, httpAction.methods...)
					if len(middlewares) > 0 {
						log.Printf("%-16s %-60s %-80s\n", strings.Join(httpAction.methods, ", "), basePath+path, strings.Join(middlewares, ", "))
					} else {
						log.Printf("%-16s %-60s\n", strings.Join(httpAction.methods, ", "), basePath+path)
					}
				}
			}
		}
	}
}

func handleHTTPAction(router *mux.Router) {
	if httpActionConfigs == nil || len(httpActionConfigs) == 0 {
		return
	}
	for _, config := range httpActionConfigs {
		subrouter := router.PathPrefix(basePath).Subrouter()
		middlewares := make([]string, 0)
		if usePanicResponse {
			subrouter.Use(mux.MiddlewareFunc(panicResponseMiddleware.Func))
			middlewares = append(middlewares, panicResponseMiddleware.Name)
		}
		for _, middleware := range config.middlewares {
			subrouter.Use(mux.MiddlewareFunc(middleware.Func))
			middlewares = append(middlewares, middleware.Name)
		}
		if useTimeLogger {
			subrouter.Use(mux.MiddlewareFunc(timeLoggerMiddleware.Func))
			middlewares = append(middlewares, timeLoggerMiddleware.Name)
		}
		handleHTTPFunc(subrouter, config.httpAction.path, config.httpAction.handlerFunc, config.httpAction.methods...)
		if len(middlewares) > 0 {
			log.Printf("%-16s %-60s %-80s\n", strings.Join(config.httpAction.methods, ", "), basePath+config.httpAction.path, strings.Join(middlewares, ", "))
		} else {
			log.Printf("%-16s %-60s\n", strings.Join(config.httpAction.methods, ", "), basePath+config.httpAction.path)
		}
	}
}

func handleHTTPFunc(subrouter *mux.Router, path string, handlerFunc http.HandlerFunc, mathods ...string) {
	subrouter.HandleFunc(path, handlerFunc).Methods(mathods...)
}

func registerRPC() {
	if rpcHandlerConfigs == nil || len(rpcHandlerConfigs) == 0 {
		return
	}
	for _, config := range rpcHandlerConfigs {
		rpcActions := []*RPCAction{}
		config.RPCHandler.Actions("all", &rpcActions)
		fmt.Println()
		rpcInfos := []string{}
		for _, action := range rpcActions {
			rpcInfos = append(rpcInfos, action.String())
		}
		if config.Name == "" {
			rpc.Register(config.RPCHandler)
			handlerType := reflect.TypeOf(config.RPCHandler)
			if handlerType.Kind() == reflect.Ptr {
				handlerType = handlerType.Elem()
			}
			name := handlerType.Name()
			log.Printf("%s: %s\n", name, strings.Join(rpcInfos, "; "))
		} else {
			rpc.RegisterName(config.Name, config.RPCHandler)
			log.Printf("%s: %s\n", config.Name, strings.Join(rpcInfos, "; "))
		}
	}
}

// ListenAndServe 监听并服务
func ListenAndServe(address string) {
	log.Println("listen on: " + address)

	router := mux.NewRouter()
	handleHTTP(router)
	handleHTTPAction(router)
	http.Handle("/", router)

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
