package gglmm

import (
	"fmt"
	"log"
	"net/rpc"
	"reflect"
	"strings"
)

// RPCAction --
type RPCAction struct {
	name     string
	request  string
	response string
}

// NewRPCAction --
func NewRPCAction(name string, request string, response string) *RPCAction {
	return &RPCAction{
		name:     name,
		request:  request,
		response: response,
	}
}

func (info RPCAction) String() string {
	return info.name + "(" + info.request + ", " + info.response + ")"
}

// RPCHandler --
type RPCHandler interface {
	Actions(cmd string, actions *[]*RPCAction) error
}

// RPCHandlerConfig --
type RPCHandlerConfig struct {
	Name       string
	RPCHandler RPCHandler
}

var rpcHandlerConfigs []*RPCHandlerConfig = nil

// RegisterRPC 注册RPCHandler
// rpcHandler 处理者
func RegisterRPC(rpcHandler RPCHandler) *RPCHandlerConfig {
	handlerType := reflect.TypeOf(rpcHandler)
	if handlerType.Kind() == reflect.Ptr {
		handlerType = handlerType.Elem()
	}
	name := handlerType.Name()
	return RegisterRPCName(name, rpcHandler)
}

// RegisterRPCName 注册RPCHandler
// name 名称
// rpcHandler 处理者
func RegisterRPCName(name string, rpcHandler RPCHandler) *RPCHandlerConfig {
	if rpcHandlerConfigs == nil {
		rpcHandlerConfigs = make([]*RPCHandlerConfig, 0)
	}
	config := &RPCHandlerConfig{
		Name:       name,
		RPCHandler: rpcHandler,
	}
	rpcHandlerConfigs = append(rpcHandlerConfigs, config)
	return config
}

func registerRPC() {
	if rpcHandlerConfigs == nil || len(rpcHandlerConfigs) == 0 {
		return
	}
	for _, config := range rpcHandlerConfigs {
		rpcActions := []*RPCAction{}
		config.RPCHandler.Actions("all", &rpcActions)
		fmt.Println()
		rpcInfos := []string{}
		for _, action := range rpcActions {
			rpcInfos = append(rpcInfos, action.String())
		}
		rpc.RegisterName(config.Name, config.RPCHandler)
		log.Printf("%s: %s\n", config.Name, strings.Join(rpcInfos, "; "))
	}
}
