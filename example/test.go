package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/rpc"
	"time"

	ws "github.com/gorilla/websocket"
	"github.com/weihongguo/gglmm"
)

func testHTTP() {
	time.Sleep(1 * time.Second)

	response, err := http.Get("http://localhost:10000/api/example/1")
	if err != nil {
		log.Println("http", err)
		return
	}
	defer response.Body.Close()

	result, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("ReadAll", err)
		return
	}

	fmt.Println()
	log.Println(string(result))
}

func testWS() {
	time.Sleep(1 * time.Second)

	client, _, err := ws.DefaultDialer.Dial("ws://localhost:10000/api/ws/example", nil)
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

func testRPC() {
	time.Sleep(3 * time.Second)

	client, err := rpc.DialHTTP("tcp", ":10000")
	if err != nil {
		log.Println("rpc", err)
		return
	}

	fmt.Println()

	idRequest := gglmm.IDRequest{
		ID: 1,
	}
	example := Example{}
	err = client.Call("ExampleRPCService.Get", idRequest, &example)
	if err != nil {
		log.Println("ExampleRPCService.Get", err)
	} else {
		log.Printf("Get: \n%+v", example)
	}

	fmt.Println()

	filterRequest := gglmm.FilterRequest{}
	filterRequest.AddFilter("id", gglmm.FilterOperateGreaterEqual, 2)
	filterRequest.AddFilter("id", gglmm.FilterOperateLessThan, 4)
	examples := make([]Example, 0)
	err = client.Call("ExampleRPCService.List", filterRequest, &examples)
	if err != nil {
		log.Println("ExampleRPCService.List", err)
	} else {
		log.Printf("List: \n%+v", examples)
	}
}
