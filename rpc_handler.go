package gglmm

// RPCActionInfo --
type RPCActionInfo struct {
	Name     string
	Request  string
	Response string
}

func (info RPCActionInfo) String() string {
	return info.Name + "(" + info.Request + ", " + info.Response + ")"
}

// RPCHandler --
type RPCHandler interface {
	Actions(cmd string, actions *[]RPCActionInfo) error
}

// RPCHandlerConfig --
type RPCHandlerConfig struct {
	RPCHandler RPCHandler
	Name       string
}
