package gglmm

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// Middleware --
type Middleware struct {
	Name string
	Func mux.MiddlewareFunc
}

// PanicResponse --
func PanicResponse() Middleware {
	return Middleware{
		Name: "PanicResponse",
		Func: func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				defer func() {
					if recover := recover(); recover != nil {
						switch panic := recover.(type) {
						case string:
							ErrorResponse(panic).
								AddData("url", r.RequestURI).
								JSON(w)
						case error:
							ErrorResponse("服务忙，请稍后再试").
								AddData("url", r.RequestURI).
								AddData("error", panic.Error()).
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

// JWTAuthentication JWT通用认证中间件
func JWTAuthentication(secrets ...string) Middleware {
	return Middleware{
		Name: fmt.Sprintf("%s%+v", "JWTAuthentication", secrets),
		Func: func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				for _, secret := range secrets {
					authorization, _, err := ParseAuthorizationToken(GetAuthorizationToken(r), secret)
					if err == nil {
						r = RequestWithAuthorization(r, authorization)
						next.ServeHTTP(w, r)
						return
					}
				}
				UnauthorizedResponse().JSON(w)
			})
		},
	}
}

// CheckPermissionFunc --
type CheckPermissionFunc func(r *http.Request) (bool, error)

// CheckPermission --
func CheckPermission(checkPermission CheckPermissionFunc) Middleware {
	return Middleware{
		Name: "CheckPermission",
		Func: func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if result, err := checkPermission(r); !result || err != nil {
					ForbiddenResponse().JSON(w)
					return
				}
				next.ServeHTTP(w, r)
			})
		},
	}
}

// TimeLogger --
func TimeLogger() Middleware {
	return Middleware{
		Name: "TimeLogger",
		Func: func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				start := time.Now().UnixNano()
				next.ServeHTTP(w, r)
				end := time.Now().UnixNano()
				log.Printf("%8.3fms %s", float64((end-start)/1000)/1000, r.RequestURI)
			})
		},
	}
}
