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

// AuthInfo --
func (user ExampleUser) AuthInfo() *auth.Info {
	return &auth.Info{
		Type: "example",
		ID:   0,
	}
}

// LoginAction --
func LoginAction(jwtExpires int64, jwtSecret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := ExampleUser{}
		authToken, _, err := auth.GenerateToken(user.AuthInfo(), jwtExpires, jwtSecret)
		if err != nil {
			gglmm.Panic(err)
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
		HTTPService: gglmm.NewHTTPService(example.Example{}),
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
