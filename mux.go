package gglmm

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// MuxVars --
func MuxVars(r *http.Request) map[string]string {
	return mux.Vars(r)
}

// MuxVar --
func MuxVar(r *http.Request, key string) (string, error) {
	vars := mux.Vars(r)
	value, ok := vars[key]
	if !ok {
		return "", errors.New("not found")
	}
	return value, nil
}

// MuxVarDefault --
func MuxVarDefault(r *http.Request, key string, defaultValue string) string {
	vars := mux.Vars(r)
	value, ok := vars[key]
	if !ok {
		return defaultValue
	}
	return value
}

// MuxVarID Mux 解释ID
func MuxVarID(r *http.Request) (int64, error) {
	value, err := MuxVar(r, "id")
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(value, 10, 64)
}
