package main

import (
	"log"
	"net/rpc"

	example "gglmm-example"

	"github.com/weihongguo/gglmm"
)

func testRPC() {

	client, err := rpc.DialHTTP("tcp", ":10000")
	if err != nil {
		log.Println("rpc", err)
		return
	}

	idRequest := gglmm.IDRequest{
		ID: 1,
	}
	one := example.Example{}
	err = client.Call("ExampleRPCService.Get", idRequest, &one)
	if err != nil {
		log.Println("ExampleRPCService.Get", err)
	} else {
		log.Printf("Get: \n%+v", one)
	}

	filterRequest := gglmm.FilterRequest{}
	filterRequest.AddFilter("id", gglmm.FilterOperateGreaterEqual, 2)
	filterRequest.AddFilter("id", gglmm.FilterOperateLessThan, 4)
	list := make([]example.Example, 0)
	err = client.Call("ExampleRPCService.List", filterRequest, &list)
	if err != nil {
		log.Println("ExampleRPCService.List", err)
	} else {
		log.Printf("List: \n%+v", list)
	}
}

func main() {
	testRPC()
}
