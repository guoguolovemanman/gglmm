package gglmm

import (
	"log"
	"net/http"

	ws "github.com/gorilla/websocket"
)

// WSMessage --
type WSMessage struct {
	Content []byte
	Close   bool
}

// NewWSMessage --
func NewWSMessage(content []byte, close bool) *WSMessage {
	return &WSMessage{
		Content: content,
		Close:   close,
	}
}

// WSHandler --
type WSHandler func(<-chan *WSMessage) <-chan *WSMessage

// WSHandlerConfig --
type WSHandlerConfig struct {
	path      string
	wsHandler WSHandler
}

var wsUpgrader *ws.Upgrader = nil
var wsHandlerConfigs []*WSHandlerConfig = nil

// HandleWS --
func HandleWS(path string, wsHandler WSHandler) *WSHandlerConfig {
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
	chanRequest := make(chan *WSMessage)

	go func() {
		defer func() {
			close(chanRequest)
			conn.Close()
		}()

		chanResponse := wsHandler(chanRequest)
		for {
			message, ok := <-chanResponse
			if message != nil || !ok {
				return
			}
			conn.WriteMessage(ws.TextMessage, message.Content)
			if message.Close {
				return
			}
		}
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("server read err:", err)
			break
		}
		log.Println("server read message")
		chanRequest <- NewWSMessage(message, false)
	}
	log.Println("server write finish")
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

func handleWS() {
	if wsHandlerConfigs == nil || len(wsHandlerConfigs) == 0 {
		return
	}
	for _, config := range wsHandlerConfigs {
		path := basePath + config.path
		log.Printf("ws %s\n", path)
		http.Handle(path, wsHandler(config.wsHandler))
	}
}
