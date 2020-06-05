package gglmm

import (
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// HTTPAction --
type HTTPAction struct {
	path        string
	handlerFunc http.HandlerFunc
	methods     []string
}

// NewHTTPAction --
func NewHTTPAction(path string, handlerFunc http.HandlerFunc, methods ...string) *HTTPAction {
	return &HTTPAction{
		path:        path,
		handlerFunc: handlerFunc,
		methods:     methods,
	}
}

// HTTPActionConfig --
type HTTPActionConfig struct {
	middlewares []Middleware
	httpAction  *HTTPAction
}

// Middleware --
func (config *HTTPActionConfig) Middleware(middlewares ...Middleware) {
	config.middlewares = middlewares
}

// HTTPHandler 提供HTTP服务接口
type HTTPHandler interface {
	Action(action Action) (*HTTPAction, error)
}

// MiddlewareAction --
type MiddlewareAction struct {
	middlewares []Middleware
	actions     []Action
}

// HTTPHandlerConfig --
type HTTPHandlerConfig struct {
	path              string
	httpHandler       HTTPHandler
	middlewareActions []MiddlewareAction
}

// Action --
func (config *HTTPHandlerConfig) Action(params ...interface{}) *HTTPHandlerConfig {
	middlewares := make([]Middleware, 0)
	actions := make([]Action, 0)
	for _, param := range params {
		if middleware, ok := param.(Middleware); ok {
			middlewares = append(middlewares, middleware)
		} else if middlewareSlice, ok := param.([]Middleware); ok {
			middlewares = append(middlewares, middlewareSlice...)
		} else if action, ok := param.(Action); ok {
			actions = append(actions, action)
		} else if actionSlice, ok := param.([]Action); ok {
			actions = append(actions, actionSlice...)
		}
	}
	if len(actions) > 0 {
		if config.middlewareActions == nil {
			config.middlewareActions = make([]MiddlewareAction, 0)
		}
		config.middlewareActions = append(config.middlewareActions, MiddlewareAction{middlewares: middlewares, actions: actions})
	}
	return config
}

var httpHandlerConfigs []*HTTPHandlerConfig = nil
var httpActionConfigs []*HTTPActionConfig = nil

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
		httpAction: &HTTPAction{
			path:        path,
			handlerFunc: handlerFunc,
			methods:     methods,
		},
	}
	httpActionConfigs = append(httpActionConfigs, config)
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
						log.Printf("http %-16s %-60s %-80s\n", strings.Join(httpAction.methods, ", "), basePath+path, strings.Join(middlewares, ", "))
					} else {
						log.Printf("http %-16s %-60s\n", strings.Join(httpAction.methods, ", "), basePath+path)
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
			log.Printf("http %-16s %-60s %-80s\n", strings.Join(config.httpAction.methods, ", "), basePath+config.httpAction.path, strings.Join(middlewares, ", "))
		} else {
			log.Printf("http %-16s %-60s\n", strings.Join(config.httpAction.methods, ", "), basePath+config.httpAction.path)
		}
	}
}

func handleHTTPFunc(subrouter *mux.Router, path string, handlerFunc http.HandlerFunc, mathods ...string) {
	subrouter.HandleFunc(path, handlerFunc).Methods(mathods...)
}
