package gglmm

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func CustomAction(w http.ResponseWriter, r *http.Request) {
	OkResponse().
		AddData("success", "success").
		JSON(w)
}

func TestGGLMM(t *testing.T) {

	HandleHTTPAction("/api/custom", CustomAction, "GET")

	router := mux.NewRouter()
	handleHTTP(router)
	handleHTTPAction(router)

	testResponse := httptest.NewRecorder()
	testRequest, _ := http.NewRequest("GET", "/api/custom", nil)

	router.ServeHTTP(testResponse, testRequest)

	response := OkResponse()
	if err := json.Unmarshal(testResponse.Body.Bytes(), response); err != nil {
		t.Fatal(err)
	}
	if response.Code != http.StatusOK {
		t.Fatal(response.Code)
	}
	success := response.Data["success"].(string)
	if success != "success" {
		t.Fatal(success)
	}
}
