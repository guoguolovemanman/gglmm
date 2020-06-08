package gglmm

import (
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

// MiddlewarePanicResponse --
func MiddlewarePanicResponse() Middleware {
	return Middleware{
		Name: "PanicResponse",
		Func: func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				defer func() {
					if recover := recover(); recover != nil {
						switch recover := recover.(type) {
						case string:
							ErrorResponse(recover).
								AddData("url", r.RequestURI).
								JSON(w)
						case ErrPanic:
							ErrorResponse(recover.message).
								AddData("url", r.RequestURI).
								AddData("file", recover.file).
								AddData("line", recover.line).
								JSON(w)
						case error:
							ErrorResponse("服务忙，请稍后再试").
								AddData("url", r.RequestURI).
								AddData("error", recover.Error()).
								JSON(w)
						default:
							ErrorResponse("服务忙，请稍后再试").
								AddData("url", r.RequestURI).
								AddData("error", "未知错误").
								JSON(w)
							log.Println(recover)
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
func MiddlewarePermissionChecker(checkPermission PermissionCheckFunc) Middleware {
	return Middleware{
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
func MiddlewareTimeLogger() Middleware {
	return Middleware{
		Name: "TimeLogger",
		Func: func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				start := time.Now().UnixNano()
				next.ServeHTTP(w, r)
				end := time.Now().UnixNano()
				log.Printf("%8.3fms %8s %s", float64((end-start)/1000)/1000, r.Method, r.RequestURI)
			})
		},
	}
}
