package gglmm

import (
	"log"

	"github.com/gorilla/mux"
)

// HTTPHandler 提供HTTP服务接口
type HTTPHandler interface {
	Action(action Action) (*HTTPAction, error)
}

// MiddlewareAction --
type MiddlewareAction struct {
	middlewares []*Middleware
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
	middlewares := make([]*Middleware, 0)
	actions := make([]Action, 0)
	for _, param := range params {
		if middleware, ok := param.(Middleware); ok {
			middlewares = append(middlewares, &middleware)
		} else if middleware, ok := param.(*Middleware); ok {
			middlewares = append(middlewares, middleware)
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

func handleHTTP(router *mux.Router) {
	if httpHandlerConfigs == nil || len(httpHandlerConfigs) == 0 {
		return
	}
	for _, config := range httpHandlerConfigs {
		subrouter := router.PathPrefix(basePath).Subrouter()
		for _, middlewareAcion := range config.middlewareActions {
			middlewares := make([]string, 0)
			if usePanicResponser {
				subrouter.Use(mux.MiddlewareFunc(middlewarePanicResponser.Func))
				middlewares = append(middlewares, middlewarePanicResponser.Name)
			}
			for _, middleware := range middlewareAcion.middlewares {
				subrouter.Use(mux.MiddlewareFunc(middleware.Func))
				middlewares = append(middlewares, middleware.Name)
			}
			if useTimeLogger {
				subrouter.Use(mux.MiddlewareFunc(middlewareTimeLogger.Func))
				middlewares = append(middlewares, middlewareTimeLogger.Name)
			}
			for _, action := range middlewareAcion.actions {
				httpAction, err := config.httpHandler.Action(action)
				if err != nil {
					log.Println(err)
				} else if httpAction.handlerFunc != nil {
					path := config.path + httpAction.path
					handleHTTPFunc(subrouter, path, httpAction.handlerFunc, httpAction.methods...)
					logHTTP(httpAction.methods, path, middlewares)
				}
			}
		}
	}
}
