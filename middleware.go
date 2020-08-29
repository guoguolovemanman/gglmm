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

// ErrFileLine --
type ErrFileLine struct {
	File    string
	Line    int
	Message string
}

func (err ErrFileLine) Error() string {
	return fmt.Sprintf("file: %s; line: %d; message: %s", err.File, err.Line, err.Message)
}

// NewErrFileLine --
func NewErrFileLine(param interface{}) *ErrFileLine {
	if _, file, line, ok := runtime.Caller(1); ok {
		switch param := param.(type) {
		case string:
			return &ErrFileLine{
				File:    file,
				Line:    line,
				Message: param,
			}
		case error:
			return &ErrFileLine{
				File:    file,
				Line:    line,
				Message: param.Error(),
			}
		default:
			return &ErrFileLine{
				File:    file,
				Line:    line,
				Message: "未知错误",
			}
		}
	} else {
		panic(&ErrFileLine{
			File:    file,
			Line:    line,
			Message: "未知错误",
		})
	}
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
						case *ErrFileLine:
							ErrorResponse(ResponseFailCode, recover.Message).
								AddData("url", r.RequestURI).
								AddData("file", recover.File).
								AddData("line", recover.Line).
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

// MiddlewareTimeLogger --
func MiddlewareTimeLogger(threshold int64) *Middleware {
	return &Middleware{
		Name: fmt.Sprintf("%s[%dms]", "TimeLogger", threshold),
		Func: func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				start := time.Now()
				defer func() {
					elapsedTime := time.Now().Sub(start)
					if elapsedTime.Nanoseconds() > threshold {
						log.Printf("%-8dms %8s %s", elapsedTime, r.Method, r.RequestURI)
					}
				}()
				next.ServeHTTP(w, r)
			})
		},
	}
}
