package main

import (
	example "gglmm-example"

	"github.com/weihongguo/gglmm"
)

// ExampleRPCService --
type ExampleRPCService struct {
	DB *gglmm.DB
}

// NewExampleRPCService --
func NewExampleRPCService() *ExampleRPCService {
	return &ExampleRPCService{
		DB: gglmm.NewDB(),
	}
}

// Actions --
func (service *ExampleRPCService) Actions(cmd string, resposne *gglmm.RPCActionsResponse) error {
	resposne.Actions = append(resposne.Actions, []*gglmm.RPCAction{
		gglmm.NewRPCAction("First", "string", "*example.Example"),
		gglmm.NewRPCAction("List", "gglmm.FilterRequest", "*[]example.Example"),
	}...)
	return nil
}

// First --
func (service *ExampleRPCService) First(idRequest *gglmm.IDRequest, example *example.Example) error {
	err := service.DB.First(example, idRequest)
	if err != nil {
		return err
	}
	return nil
}

// List --
func (service *ExampleRPCService) List(filterRequest *gglmm.FilterRequest, examples *[]example.Example) error {
	service.DB.List(examples, filterRequest)
	return nil
}
