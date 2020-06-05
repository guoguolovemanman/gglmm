package main

import "github.com/weihongguo/gglmm"

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
func (service *ExampleRPCService) Actions(cmd string, actions *[]*gglmm.RPCAction) error {
	*actions = append(*actions, []*gglmm.RPCAction{
		gglmm.NewRPCAction("Get", "string", "*Test"),
		gglmm.NewRPCAction("List", "gglmm.FilterRequest", "*[]Test"),
	}...)
	return nil
}

// Get --
func (service *ExampleRPCService) Get(idRequest gglmm.IDRequest, example *Example) error {
	err := service.gormDB.Get(example, idRequest)
	if err != nil {
		return err
	}
	return nil
}

// List --
func (service *ExampleRPCService) List(filterRequest gglmm.FilterRequest, examples *[]Example) error {
	service.gormDB.List(examples, filterRequest)
	return nil
}
