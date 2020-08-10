package main

import (
	"net/http"

	example "gglmm-example"

	"github.com/weihongguo/gglmm"
	auth "github.com/weihongguo/gglmm-auth"
)

// ExampleUser --
type ExampleUser struct {
}

// Login --
func (user ExampleUser) Login() *auth.Subject {
	return &auth.Subject{
		UserType: "example",
		UserID:   0,
	}
}

// LoginAction --
func LoginAction(jwtExpires int64, jwtSecret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := ExampleUser{}
		authToken, _, err := auth.GenerateToken(user.Login(), jwtExpires, jwtSecret)
		if err != nil {
			gglmm.FailResponse(gglmm.NewErrFileLine(err)).JSON(w)
			return
		}
		gglmm.OkResponse().
			AddData("authToken", authToken).
			JSON(w)
	}
}

// ExampleService --
type ExampleService struct {
	*gglmm.HTTPService
}

// NewExampleService --
func NewExampleService() *ExampleService {
	return &ExampleService{
		HTTPService: gglmm.NewHTTPService(example.Example{}, [...]string{"example", "examples"}),
	}
}

// ExampleAction --
func (service *ExampleService) ExampleAction(w http.ResponseWriter, r *http.Request) {
	gglmm.OkResponse().JSON(w)
}

// ExampleAction --
func ExampleAction(w http.ResponseWriter, r *http.Request) {
	gglmm.OkResponse().JSON(w)
}
