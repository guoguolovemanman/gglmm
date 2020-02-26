package gglmm

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type TestService struct{}

func (service *TestService) CustomActions() ([]*HTTPAction, error) {
	actions := []*HTTPAction{
		&HTTPAction{"/test", func(w http.ResponseWriter, r *http.Request) {
			NewSuccessResponse().
				AddData("success", "success").
				WriteJSON(w)
		}, "GET"},
	}
	return actions, nil
}

func (service *TestService) RESTAction(action RESTAction) (*HTTPAction, error) {
	return nil, nil
}

func TestHTTP(t *testing.T) {

	RegisterHTTPHandler(&TestService{}, "/api")

	router := GenerateHttpRouter()

	testResponse := httptest.NewRecorder()
	testRequest, _ := http.NewRequest("GET", "/api/test", nil)

	router.ServeHTTP(testResponse, testRequest)

	response := NewResponse()
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
