package gglmm

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/gorilla/mux"
)

// Middleware --
type Middleware struct {
	Name string
	Func mux.MiddlewareFunc
}

// ErrPanic --
type ErrPanic struct {
	file    string
	line    int
	message string
}

// Panic --
func Panic(param interface{}) {
	errPanic := ErrPanic{}
	switch err := param.(type) {
	case string:
		errPanic.message = err
	case error:
		errPanic.message = err.Error()
	default:
		errPanic.message = "服务忙，请稍后再试"
	}
	if _, file, line, ok := runtime.Caller(1); ok {
		errPanic.file = file
		errPanic.line = line
	}
	panic(errPanic)
}

// MiddlewarePanicResponser --
func MiddlewarePanicResponser() *Middleware {
	return &Middleware{
		Name: "PanicResponser",
		Func: func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				defer func() {
					if recover := recover(); recover != nil {
						switch recover := recover.(type) {
						case string:
							ErrorResponse(ResponseFailCode, recover).
								AddData("url", r.RequestURI).
								JSON(w)
						case ErrPanic:
							ErrorResponse(ResponseFailCode, recover.message).
								AddData("url", r.RequestURI).
								AddData("file", recover.file).
								AddData("line", recover.line).
								JSON(w)
						case error:
							ErrorResponse(ResponseFailCode, "服务忙，请稍后再试").
								AddData("url", r.RequestURI).
								AddData("error", recover.Error()).
								JSON(w)
						default:
							ErrorResponse(ResponseFailCode, "服务忙，请稍后再试").
								AddData("url", r.RequestURI).
								AddData("error", "未知错误").
								JSON(w)
						}
					}
				}()
				next.ServeHTTP(w, r)
			})
		},
	}
}

// PermissionCheckFunc --
type PermissionCheckFunc func(r *http.Request) (bool, error)

// MiddlewarePermissionChecker --
func MiddlewarePermissionChecker(checkPermission PermissionCheckFunc) *Middleware {
	return &Middleware{
		Name: "PermissionChecker",
		Func: func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if result, err := checkPermission(r); !result || err != nil {
					log.Println(err)
					ForbiddenResponse().JSON(w)
					return
				}
				next.ServeHTTP(w, r)
			})
		},
	}
}

// MiddlewareTimeLogger --
func MiddlewareTimeLogger(threshold int64) *Middleware {
	return &Middleware{
		Name: fmt.Sprintf("%s[%dms]", "TimeLogger", threshold),
		Func: func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				start := time.Now().UnixNano()
				defer func() {
					end := time.Now().UnixNano()
					elapsedTime := (end - start) / 1000 / 1000
					if elapsedTime > threshold {
						log.Printf("%-8dms %8s %s", elapsedTime, r.Method, r.RequestURI)
					}
				}()
				next.ServeHTTP(w, r)
			})
		},
	}
}
