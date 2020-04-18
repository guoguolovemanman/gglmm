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
var httpHandlerConfigs []*HTTPHandlerConfig = nil
var rpcHandlerConfigs []*RPCHandlerConfig = nil

// BasePath 注册基础路径
func BasePath(path string) {
	basePath = path
}

// HandleHTTP 注册HTTP请求处理者
// httpHandler 处理者
// path 路径
func HandleHTTP(httpHandler HTTPHandler, params ...interface{}) *HTTPHandlerConfig {
	if httpHandlerConfigs == nil {
		httpHandlerConfigs = make([]*HTTPHandlerConfig, 0)
	}
	config := &HTTPHandlerConfig{
		HTTPHandler: httpHandler,
	}
	if params != nil {
		if len(params) > 0 {
			if path, ok := params[0].(string); ok {
				config.Path = path
			}
		}
	}
	httpHandlerConfigs = append(httpHandlerConfigs, config)
	return config
}

// RegisterRPC 注册RPC请求处理者
// rpcHandler 处理者
// name 名称
func RegisterRPC(rpcHandler RPCHandler, params ...interface{}) *RPCHandlerConfig {
	if rpcHandlerConfigs == nil {
		rpcHandlerConfigs = make([]*RPCHandlerConfig, 0)
	}
	config := &RPCHandlerConfig{
		RPCHandler: rpcHandler,
	}
	if params != nil {
		if len(params) > 0 {
			if name, ok := params[0].(string); ok {
				config.Name = name
			}
		}
	}
	rpcHandlerConfigs = append(rpcHandlerConfigs, config)
	return config
}

func handleHTTP() *mux.Router {
	if httpHandlerConfigs == nil || len(httpHandlerConfigs) == 0 {
		return nil
	}

	router := mux.NewRouter()
	for _, config := range httpHandlerConfigs {
		subrouter := router.PathPrefix(basePath).Subrouter()
		var middlewares string
		for _, middleware := range config.Middlewares {
			subrouter.Use(mux.MiddlewareFunc(middleware.Func))
			if middlewares == "" {
				middlewares = middleware.Name
			} else {
				middlewares += "|" + middleware.Name
			}
		}
		fmt.Println()
		httpActions, err := config.HTTPHandler.CustomActions()
		if err != nil {
			log.Println(err)
		} else {
			if httpActions != nil {
				for _, httpAction := range httpActions {
					handleHTTPAction(subrouter, middlewares, config, httpAction)
				}
			}
		}
		for _, action := range config.Actions {
			httpAction, err := config.HTTPHandler.Action(action)
			if err != nil {
				log.Println(err)
			} else if httpAction.HandlerFunc != nil {
				handleHTTPAction(subrouter, middlewares, config, httpAction)
			}
		}
	}
	return router
}

func handleHTTPAction(subrouter *mux.Router, middlewares string, config *HTTPHandlerConfig, httpAction *HTTPAction) {
	if httpAction.HandlerFunc == nil {
		return
	}
	path := config.Path + httpAction.Path
	subrouter.HandleFunc(path, httpAction.HandlerFunc).Methods(httpAction.Method)
	if middlewares != "" {
		log.Printf("%-8s %-60s %-40s\n", httpAction.Method, basePath+path, middlewares)
	} else {
		log.Printf("%-8s %-60s\n", httpAction.Method, basePath+path)
	}
}

func registerRPC() {
	for _, config := range rpcHandlerConfigs {
		if config.Name == "" {
			rpc.Register(config.RPCHandler)
		} else {
			rpc.RegisterName(config.Name, config.RPCHandler)
		}
		rpcActionInfos := []RPCActionInfo{}
		config.RPCHandler.Actions("all", &rpcActionInfos)
		fmt.Println()
		rpcInfos := []string{}
		for _, info := range rpcActionInfos {
			rpcInfos = append(rpcInfos, info.String())
		}
		if config.Name == "" {
			handlerType := reflect.TypeOf(config.RPCHandler)
			if handlerType.Kind() == reflect.Ptr {
				handlerType = handlerType.Elem()
			}
			name := handlerType.Name()
			log.Printf("%s: %s\n", name, strings.Join(rpcInfos, "; "))
		} else {
			log.Printf("%s: %s\n", config.Name, strings.Join(rpcInfos, "; "))
		}
	}
}

// ListenAndServe 监听并服务
func ListenAndServe(address string) {
	log.Println("listen on: " + address)
	if httpHandlerConfigs != nil && len(httpHandlerConfigs) >= 0 {
		router := handleHTTP()
		http.Handle("/", router)
	}
	if rpcHandlerConfigs != nil && len(rpcHandlerConfigs) >= 0 {
		registerRPC()
		rpc.HandleHTTP()
	}
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
