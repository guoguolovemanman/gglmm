package gglmm

import "github.com/gorilla/mux"

// Middleware --
type Middleware struct {
	Name string
	Func mux.MiddlewareFunc
}
