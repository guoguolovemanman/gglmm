package main

import (
	example "gglmm-example"

	"github.com/weihongguo/gglmm"
)

// ExampleRPCService --
type ExampleRPCService struct {
	gormDB *gglmm.GormDB
}

// NewExampleRPCService --
func NewExampleRPCService() *ExampleRPCService {
	return &ExampleRPCService{
		gormDB: gglmm.DefaultGormDB(),
	}
}

// Actions --
func (service *ExampleRPCService) Actions(cmd string, resposne *gglmm.RPCActionsResponse) error {
	resposne.Actions = append(resposne.Actions, []*gglmm.RPCAction{
		gglmm.NewRPCAction("Get", "string", "*example.Example"),
		gglmm.NewRPCAction("List", "gglmm.FilterRequest", "*[]example.Example"),
	}...)
	return nil
}

// Get --
func (service *ExampleRPCService) Get(idRequest *gglmm.IDRequest, example *example.Example) error {
	err := service.gormDB.Get(example, idRequest)
	if err != nil {
		return err
	}
	return nil
}

// List --
func (service *ExampleRPCService) List(filterRequest *gglmm.FilterRequest, examples *[]example.Example) error {
	service.gormDB.List(examples, filterRequest)
	return nil
}
