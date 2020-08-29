package gglmm

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

// ErrPathVar --
var ErrPathVar = errors.New("路径参数错误")

// PathVars --
func PathVars(r *http.Request) map[string]string {
	return mux.Vars(r)
}

// PathVar --
func PathVar(r *http.Request, key string) (string, error) {
	vars := PathVars(r)
	value, ok := vars[key]
	if !ok {
		return "", ErrPathVar
	}
	return value, nil
}

// PathVarID Mux 解释ID
func PathVarID(r *http.Request) (uint64, error) {
	value, err := PathVar(r, "id")
	if err != nil {
		return 0, err
	}
	result, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func handleHTTPFunc(subrouter *mux.Router, path string, handlerFunc http.HandlerFunc, mathods ...string) {
	subrouter.HandleFunc(path, handlerFunc).Methods(mathods...)
}

func logHTTP(methods []string, path string, middlewares []string) {
	if len(middlewares) > 0 {
		log.Printf("[http] [%-16s] %-60s %-80s\n", strings.Join(methods, ", "), basePath+path, strings.Join(middlewares, ", "))
	} else {
		log.Printf("[http] [%-16s] %-60s\n", strings.Join(methods, ", "), basePath+path)
	}
}
