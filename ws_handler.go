package gglmm

import (
	"log"
	"net/http"

	ws "github.com/gorilla/websocket"
)

// WSHandler --
type WSHandler func(conn *ws.Conn, messageType int, message []byte) bool

// WSHandlerConfig --
type WSHandlerConfig struct {
	path      string
	wsHandler WSHandler
}

var wsUpgrader *ws.Upgrader = nil
var wsHandlerConfigs []*WSHandlerConfig = nil

// RegisterWS --
func RegisterWS(path string, wsHandler WSHandler) *WSHandlerConfig {
	if wsUpgrader == nil {
		wsUpgrader = &ws.Upgrader{}
	}
	if wsHandlerConfigs == nil {
		wsHandlerConfigs = make([]*WSHandlerConfig, 0)
	}
	config := &WSHandlerConfig{
		path:      path,
		wsHandler: wsHandler,
	}
	wsHandlerConfigs = append(wsHandlerConfigs, config)
	return config
}

func dealMessage(conn *ws.Conn, wsHandler WSHandler) {
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}
		if wsHandler(conn, messageType, message) {
			break
		}
	}
}

func wsHandler(wsHandler WSHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := wsUpgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		go dealMessage(conn, wsHandler)
	}
}

func registerWS() {
	if wsHandlerConfigs == nil || len(wsHandlerConfigs) == 0 {
		return
	}
	for _, config := range wsHandlerConfigs {
		log.Printf("websocket %s\n", config.path)
		http.Handle(config.path, wsHandler(config.wsHandler))
	}
}
