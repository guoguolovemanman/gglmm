package gglmm

import (
	"fmt"
	"log"
	"net/http"
	"net/rpc"

	"github.com/gorilla/mux"
)

var basePath string
var httpHandlerConfigs []*HTTPHandlerConfig
var rpcHandlerConfigs []*RPCHandlerConfig

// RegisterBasePath --
func RegisterBasePath(path string) {
	basePath = path
}

// RegisterHTTPHandler --
func RegisterHTTPHandler(httpHandler HTTPHandler, path string) *HTTPHandlerConfig {
	if httpHandlerConfigs == nil {
		httpHandlerConfigs = make([]*HTTPHandlerConfig, 0)
	}
	config := &HTTPHandlerConfig{
		HTTPHandler: httpHandler,
		Path:        path,
	}
	httpHandlerConfigs = append(httpHandlerConfigs, config)
	return config
}

// RegisterRPCHandler --
func RegisterRPCHandler(rpcHandler RPCHandler, name string) *RPCHandlerConfig {
	if rpcHandlerConfigs == nil {
		rpcHandlerConfigs = make([]*RPCHandlerConfig, 0)
	}
	config := &RPCHandlerConfig{
		RPCHandler: rpcHandler,
		Name:       name,
	}
	rpcHandlerConfigs = append(rpcHandlerConfigs, config)
	return config
}

// GenerateRouter --
func GenerateRouter() *mux.Router {
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
					setupAction(subrouter, middlewares, config, httpAction)
				}
			}
		}

		for _, restAction := range config.RESTActions {
			httpAction, err := config.HTTPHandler.RESTAction(restAction)
			if err != nil {
				log.Println(err)
			} else if httpAction.HandlerFunc != nil {
				setupAction(subrouter, middlewares, config, httpAction)
			}
		}
	}
	return router
}

func setupAction(subrouter *mux.Router, middlewares string, config *HTTPHandlerConfig, httpAction *HTTPAction) {
	if httpAction.HandlerFunc == nil {
		return
	}
	path := config.Path + httpAction.Path
	subrouter.HandleFunc(path, httpAction.HandlerFunc).Methods(httpAction.Method)
	if middlewares != "" {
		log.Printf("%-8s %-60s %-60s\n", httpAction.Method, basePath+path, middlewares)
	} else {
		log.Printf("%-8s %-60s\n", httpAction.Method, basePath+path)
	}
}

// ListenAndServe 监听并服务
func ListenAndServe(address string) {
	fmt.Println()
	log.Println("listen on: " + address)

	if httpHandlerConfigs != nil && len(httpHandlerConfigs) >= 0 {
		router := GenerateRouter()
		http.Handle("/", router)
	}

	if rpcHandlerConfigs != nil && len(rpcHandlerConfigs) >= 0 {
		info := RPCInfo{}
		for _, config := range rpcHandlerConfigs {
			rpc.RegisterName(config.Name, config.RPCHandler)

			config.RPCHandler.Info("all", &info)

			fmt.Println()
			log.Printf("Name: %s; %s\n", config.Name, info.Info)
		}
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
		log.Printf("%+v\n", config)
		log.Fatal("APIConfig invalid")
	}
	ListenAndServe(config.Address)
}
