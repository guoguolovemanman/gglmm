package gglmm

import (
	"encoding/json"
	"testing"
)

func TestResponse(t *testing.T) {
	response := OkResponse().AddData("test", "test").AddData("haha", "haha")
	jsonStr, err := json.Marshal(response)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(jsonStr))
}
