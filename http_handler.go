package gglmm

import "net/http"

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
