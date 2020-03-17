package gglmm

// RPCHandler --
type RPCHandler interface {
	Actions(cmd string, actions *[]string) error
}

// RPCHandlerConfig --
type RPCHandlerConfig struct {
	RPCHandler RPCHandler
	Name       string
}
