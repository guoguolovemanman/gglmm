package gglmm

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type TestService struct{}

func (service *TestService) CustomActions() ([]*HTTPAction, error) {
	actions := []*HTTPAction{}
	action := &HTTPAction{"/test", func(w http.ResponseWriter, r *http.Request) {
		OkResponse().
			AddData("success", "success").
			JSON(w)
	}, "GET"}
	actions = append(actions, action)
	return actions, nil
}

func (service *TestService) Action(action string) (*HTTPAction, error) {
	return nil, nil
}

func TestHTTP(t *testing.T) {

	HandleHTTP(&TestService{}, "/api")

	router := handleHTTP()

	testResponse := httptest.NewRecorder()
	testRequest, _ := http.NewRequest("GET", "/api/test", nil)

	router.ServeHTTP(testResponse, testRequest)

	response := OkResponse()
	if err := json.Unmarshal(testResponse.Body.Bytes(), response); err != nil {
		t.Fail()
	}
	if response.Code != http.StatusOK {
		t.Fail()
	}
	success := response.Data["success"].(string)
	if success != "success" {
		t.Fail()
	}
}
