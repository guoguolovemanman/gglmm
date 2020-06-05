package main

import (
	"log"
	"time"

	"github.com/weihongguo/gglmm"
)

// OnceWSHandler --
func OnceWSHandler(chanResponse chan<- *gglmm.WSMessage, chanRequest <-chan *gglmm.WSMessage) {
	chanOver := time.After(5 * time.Second)
	for {
		select {
		case message, ok := <-chanRequest:
			if !ok {
				return
			}
			// 逻辑
			if message.Content != nil && string(message.Content) == "have" {
				log.Println("server handler success")
				chanResponse <- gglmm.NewWSMessage([]byte("success"), true)
				return
			}
			if message.Over {
				return
			}
		case <-chanOver:
			// 超时
			log.Println("server handler timout")
			chanResponse <- gglmm.NewWSMessage([]byte("timout"), true)
			return
		}
	}
}

// EchoWSHandler --
func EchoWSHandler(chanResponse chan<- *gglmm.WSMessage, chanRequest <-chan *gglmm.WSMessage) {
	for {
		select {
		case message, ok := <-chanRequest:
			if !ok {
				return
			}
			if message.Content != nil {
				chanResponse <- gglmm.NewWSMessage(message.Content, false)
			}
			if message.Over {
				return
			}
		}
	}
}
