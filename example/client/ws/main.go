package main

import (
	"log"
	"time"

	ws "github.com/gorilla/websocket"
)

func testEchoWS() {

	client, _, err := ws.DefaultDialer.Dial("ws://localhost:10000/api/ws/echo", nil)
	if err != nil {
		log.Println("#", err)
		return
	}
	defer client.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			messageType, message, err := client.ReadMessage()
			if err != nil {
				log.Println("client read err:", err)
				return
			}
			log.Println("client read: ", messageType, string(message))
		}
	}()

	ticker := time.NewTicker(time.Second)
	for {
		select {
		case t := <-ticker.C:
			messageType := ws.TextMessage
			message := []byte(t.String())
			err := client.WriteMessage(messageType, message)
			if err != nil {
				log.Println("client write err:", err)
				return
			}
			log.Println("client write: ", messageType, string(message))
		case <-done:
			log.Println("client write finish")
			return
		}
	}
}

func testOnceWS() {

	client, _, err := ws.DefaultDialer.Dial("ws://localhost:10000/api/ws/once", nil)
	if err != nil {
		log.Println("#", err)
		return
	}
	defer client.Close()

	err = client.WriteMessage(ws.TextMessage, []byte("have"))
	if err != nil {
		log.Println("client write err:", err)
		return
	}

	messageType, message, err := client.ReadMessage()
	if err != nil {
		log.Println("client read err:", err)
		return
	}
	log.Println("client read: ", messageType, string(message))
}

func main() {
	// testEchoWS()
	testOnceWS()
}
