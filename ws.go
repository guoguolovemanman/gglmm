package gglmm

import (
	"log"
	"net/http"

	ws "github.com/gorilla/websocket"
)

// WSMessage --
type WSMessage struct {
	Content []byte
	Over    bool
}

// NewWSMessage --
func NewWSMessage(content []byte, over bool) *WSMessage {
	return &WSMessage{
		Content: content,
		Over:    over,
	}
}

// WSHandler --
type WSHandler func(chanResponse chan<- *WSMessage, chanRequest <-chan *WSMessage)

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

func messageTransfer(conn *ws.Conn, wsHandler WSHandler) {
	chanRequest := make(chan *WSMessage)
	chanResponse := make(chan *WSMessage)

	defer func() {
		log.Println("server messageTransfer close channel")
		close(chanRequest)
		close(chanResponse)
	}()

	go func() {
		wsHandler(chanResponse, chanRequest)
		log.Println("server wsHandler finish")
	}()

	go func() {
		for {
			message, ok := <-chanResponse
			if !ok {
				conn.Close()
				return
			}
			if message.Content != nil {
				log.Println("server send message", string(message.Content))
				conn.WriteMessage(ws.TextMessage, message.Content)
			}
			if message.Over {
				conn.Close()
				return
			}
		}
	}()

	for {
		_, content, err := conn.ReadMessage()
		if err != nil {
			log.Println("server read err:", err)
			break
		}
		log.Println("server receive message", string(content))
		chanRequest <- NewWSMessage(content, false)
	}
	log.Println("server messageTransfer finish")
}

func wsHandler(wsHandler WSHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := wsUpgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("server new conn")
		go messageTransfer(conn, wsHandler)
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
