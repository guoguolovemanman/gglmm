package gglmm

import (
	"errors"
	"net/http"
	"strconv"

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
