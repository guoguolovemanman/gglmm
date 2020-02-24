package gglmm

// RPCInfo --
type RPCInfo struct {
	Info string
}

func (info RPCInfo) String() string {
	return info.Info
}

// RPCHandler --
type RPCHandler interface {
	Info(cmd string, info *RPCInfo) error
}

// RPCHandlerConfig --
type RPCHandlerConfig struct {
	Name       string
	RPCHandler RPCHandler
}
