package main

import (
	"bytes"
	"encoding/json"
	example "gglmm-example"
	"io/ioutil"
	"log"
	"net/http"

	auth "github.com/weihongguo/gglmm-auth"
)

func testHTTP() {
	// 没有登录态
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
	log.Print(string(result))

	client := &http.Client{}

	// 取登录态
	loginRequest := auth.LoginRequest{
		UserName: "example",
		Password: "example",
	}
	body, err := json.Marshal(loginRequest)
	request, err := http.NewRequest("POST", "http://localhost:10000/api/login", bytes.NewReader(body))
	response, err = client.Do(request)
	if err != nil {
		log.Println("http", err)
		return
	}
	defer response.Body.Close()

	result, err = ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("ReadAll", err)
		return
	}

	loginRespone := make(map[string]interface{})
	err = json.Unmarshal(result, &loginRespone)
	if err != nil {
		log.Println("json", err)
		return
	}
	log.Println("loginRespone", loginRespone)
	data, ok := loginRespone["data"]
	if !ok {
		log.Println("login fail")
		return
	}
	loginData, ok := data.(map[string]interface{})
	if !ok {
		log.Println("login fail")
		return
	}
	token, ok := loginData["authToken"].(string)
	if !ok {
		log.Println("login fail")
		return
	}
	log.Println(token)

	// 带上登录态请求

	// create
	example := example.Example{
		IntValue:    1,
		StringValue: "1",
		FloatValue:  1.1,
	}
	body, err = json.Marshal(example)
	request, err = http.NewRequest("POST", "http://localhost:10000/api/example", bytes.NewReader(body))
	request.Header.Add("Authorization", "Bearer "+token)
	response, err = client.Do(request)
	if err != nil {
		log.Println("http", err)
		return
	}
	defer response.Body.Close()

	result, err = ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("ReadAll", err)
		return
	}
	log.Print(string(result))

	// get
	request, err = http.NewRequest("GET", "http://localhost:10000/api/example/1", nil)
	if err != nil {
		log.Println("http", err)
		return
	}
	request.Header.Add("Authorization", "Bearer "+token)
	response, err = client.Do(request)
	if err != nil {
		log.Println("http", err)
		return
	}
	defer response.Body.Close()

	result, err = ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("ReadAll", err)
		return
	}
	log.Print(string(result))
}

func main() {
	testHTTP()
}
