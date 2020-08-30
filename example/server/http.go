package main

import (
	"log"
	"net/http"

	example "gglmm-example"

	"github.com/weihongguo/gglmm"
	auth "github.com/weihongguo/gglmm-auth"
)

// ExampleUser --
type ExampleUser struct {
	gglmm.Model
}

func (user *ExampleUser) authInfo() (*auth.Info, error) {
	return &auth.Info{
		Subject: &auth.Subject{
			Project:  "example",
			UserType: "example",
			UserID:   1,
		},
		Nickname:  "example",
		AvatarURL: "example",
	}, nil
}

// Login --
func (user *ExampleUser) Login(request *auth.LoginRequest) (*auth.Info, error) {
	log.Println("ExampleUser.Login")
	return user.authInfo()
}

// Info --
func (user *ExampleUser) Info(request *gglmm.IDRequest) (*auth.Info, error) {
	log.Println("ExampleUser.Info")
	return user.authInfo()
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
