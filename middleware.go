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
type CheckPermissionFunc func(r *http.Request) error

// PermissionChecker --
func PermissionChecker(checkPermission CheckPermissionFunc) Middleware {
	return Middleware{
		Name: "PermissionChecker",
		Func: func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if err := checkPermission(r); err != nil {
					ForbiddenResponse().JSON(w)
					return
				}
				next.ServeHTTP(w, r)
			})
		},
	}
}
