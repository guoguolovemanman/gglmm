package gglmm

import "net/http"

// HTTPAction --
type HTTPAction struct {
	Path        string
	HandlerFunc http.HandlerFunc
	Method      string
}

// NewHTTPAction --
func NewHTTPAction(path string, handlerFunc http.HandlerFunc, method string) *HTTPAction {
	return &HTTPAction{
		Path:        path,
		HandlerFunc: handlerFunc,
		Method:      method,
	}
}

// HTTPHandler 提供HTTP服务接口
type HTTPHandler interface {
	CustomActions() ([]*HTTPAction, error)
	Action(action string) (*HTTPAction, error)
}

// HTTPHandlerConfig --
type HTTPHandlerConfig struct {
	HTTPHandler HTTPHandler
	Path        string
	Actions     []string
	Middlewares []Middleware
}

// Action --
func (config *HTTPHandlerConfig) Action(actions ...string) *HTTPHandlerConfig {
	if config.Actions == nil {
		config.Actions = make([]string, 0)
	}
	config.Actions = append(config.Actions, actions...)
	return config
}

// Middleware --
func (config *HTTPHandlerConfig) Middleware(middlewares ...Middleware) *HTTPHandlerConfig {
	if config.Middlewares == nil {
		config.Middlewares = make([]Middleware, 0)
	}
	config.Middlewares = append(config.Middlewares, middlewares...)
	return config
}
