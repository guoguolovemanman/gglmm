package main

import (
	"fmt"
	"log"
	"net/rpc"

	"github.com/weihongguo/gglmm"
)

// Example --
type Example struct {
	gglmm.Model
	IntValue    int     `json:"intValue"`
	FloatValue  float64 `json:"floatValue"`
	StringValue string  `json:"stringValue"`
}

func testRPC() {

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

func main() {
	testRPC()
}
