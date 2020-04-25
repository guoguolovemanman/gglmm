package gglmm

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
