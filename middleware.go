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

// JWTAuthMiddleware JWT通用认证中间件
func JWTAuthMiddleware(secrets []string) Middleware {
	return Middleware{
		Name: fmt.Sprintf("%s%+v", "JWTAuth", secrets),
		Func: func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				for _, secret := range secrets {
					authInfo, _, err := ParseAuthToken(GetAuthToken(r), secret)
					if err == nil {
						r = RequestWithAuthInfo(r, authInfo)
						next.ServeHTTP(w, r)
						return
					}
				}
				UnauthorizedResponse().JSON(w)
			})
		},
	}
}
