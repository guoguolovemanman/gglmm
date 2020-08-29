package gglmm

import (
	"net/http"

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
	middlewares []*Middleware
	httpAction  *HTTPAction
}

// Middleware --
func (config *HTTPActionConfig) Middleware(middlewares ...*Middleware) {
	config.middlewares = middlewares
}

var httpActionConfigs []*HTTPActionConfig = nil

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

func handleHTTPAction(router *mux.Router) {
	if httpActionConfigs == nil || len(httpActionConfigs) == 0 {
		return
	}
	for _, config := range httpActionConfigs {
		subrouter := router.PathPrefix(basePath).Subrouter()
		middlewares := make([]string, 0)
		if usePanicResponser {
			subrouter.Use(mux.MiddlewareFunc(MiddlewarePanicResponser().Func))
			middlewares = append(middlewares, middlewarePanicResponser.Name)
		}
		for _, middleware := range config.middlewares {
			subrouter.Use(mux.MiddlewareFunc(middleware.Func))
			middlewares = append(middlewares, middleware.Name)
		}
		if useTimeLogger {
			subrouter.Use(mux.MiddlewareFunc(middlewareTimeLogger.Func))
			middlewares = append(middlewares, middlewareTimeLogger.Name)
		}
		handleHTTPFunc(subrouter, config.httpAction.path, config.httpAction.handlerFunc, config.httpAction.methods...)
		logHTTP(config.httpAction.methods, config.httpAction.path, middlewares)
	}
}
