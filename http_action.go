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
	httpAction  HTTPAction
}

// Middleware --
func (config *HTTPActionConfig) Middleware(middlewares ...Middleware) {
	config.middlewares = middlewares
}
