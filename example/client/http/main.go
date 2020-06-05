package main

import (
	"io/ioutil"
	"log"
	"net/http"
)

func testHTTP() {

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

	log.Println(string(result))
}

func main() {
	testHTTP()
}
