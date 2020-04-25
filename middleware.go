package gglmm

import (
	"fmt"
	"net/http"

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
							InternalErrorResponse(panic).JSON(w)
						case error:
							InternalErrorResponse(panic.Error()).JSON(w)
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
