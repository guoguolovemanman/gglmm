package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
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

	// 取登录态
	response, err = http.Get("http://localhost:10000/api/login")
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

	loginResult := make(map[string]interface{})
	err = json.Unmarshal(result, &loginResult)
	if err != nil {
		log.Println("json", err)
		return
	}
	data, ok := loginResult["data"]
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
	client := &http.Client{}
	request, err := http.NewRequest("GET", "http://localhost:10000/api/example/1", nil)
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
