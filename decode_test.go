package gglmm

import (
	"bytes"
	"encoding/json"
	"net/http"
	"reflect"
	"testing"
)

func TestDecodeIDRequest(t *testing.T) {
	url := "/test?id=1&preloads=a,b"
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}
	resultRequest, err := DecodeIDRequest(request)
	if err != nil {
		t.Fatal(err)
	}
	if resultRequest.ID != 1 {
		t.Fatal(resultRequest)
	}
	if len(resultRequest.Preloads) != 2 {
		t.Fatal(resultRequest)
	}
	if resultRequest.Preloads[0] != "a" || resultRequest.Preloads[1] != "b" {
		t.Fatal(resultRequest)
	}
}

func TestDecodeFilterRequest(t *testing.T) {
	filterRequest := FilterRequest{
		Filters: []Filter{
			{Field: "A", Operate: FilterOperateEqual, Value: "B"},
		},
		Order: "id",
	}
	query, _ := json.Marshal(filterRequest)
	url := "/test"
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(query))
	if err != nil {
		t.Fatal(err)
	}
	resultRequest, err := DecodeFilterRequest(request)
	if err != nil {
		t.Fatal(err)
	}
	if resultRequest.Order != "id" {
		t.Fatal(resultRequest)
	}
	filter := resultRequest.Filters[0]
	if filter.Field != "A" || filter.Operate != FilterOperateEqual || filter.Value != "B" {
		t.Fatal(resultRequest)
	}
}

func TestDecodePageRequest(t *testing.T) {
	pageRequest := PageRequest{
		FilterRequest: &FilterRequest{
			Filters: []Filter{
				{Field: "A", Operate: FilterOperateEqual, Value: "B"},
			},
			Order: "id",
		},
		Pagination: Pagination{
			PageSize:  0,
			PageIndex: 1,
			Total:     2,
		},
	}
	query, _ := json.Marshal(pageRequest)
	url := "/test"
	request, err := http.NewRequest("POST", url, bytes.NewReader(query))
	if err != nil {
		t.Fatal(err)
	}
	resultRequest, err := DecodePageRequest(request)
	if err != nil {
		t.Fatal(err)
	}
	filterRequest := resultRequest.FilterRequest
	if filterRequest.Order != "id" {
		t.Fatal(resultRequest)
	}
	filter := filterRequest.Filters[0]
	if filter.Field != "A" || filter.Operate != FilterOperateEqual || filter.Value != "B" {
		t.Fatal(resultRequest)
	}
	pagination := resultRequest.Pagination
	if pagination.PageSize != 0 || pagination.PageIndex != 1 || pagination.Total != 2 {
		t.Fatal(resultRequest)
	}
}

type TestModel struct {
	Key   string
	Value string
}

func TestDecodeModel(t *testing.T) {
	test := TestModel{Key: "key", Value: "value"}
	query, _ := json.Marshal(test)
	request, err := http.NewRequest("POST", "", bytes.NewReader(query))
	if err != nil {
		t.Fatal(err)
	}
	result, err := DecodeModel(request, reflect.TypeOf(test))
	if err != nil {
		t.Fatal(err)
	}
	check := result.(TestModel)
	if check.Key != test.Key || check.Value != test.Value {
		t.Fatal(check)
	}
}

func TestDecodeModelPtr(t *testing.T) {
	test := TestModel{Key: "key", Value: "value"}
	query, _ := json.Marshal(test)
	request, err := http.NewRequest("POST", "", bytes.NewReader(query))
	if err != nil {
		t.Fatal(err)
	}
	result, err := DecodeModelPtr(request, reflect.TypeOf(test))
	if err != nil {
		t.Fatal(err)
	}
	check := result.(*TestModel)
	if check.Key != test.Key || check.Value != test.Value {
		t.Fatal(check)
	}
}

func TestDecodeModelSlice(t *testing.T) {
	tests := []TestModel{
		{Key: "key", Value: "value"},
	}
	query, _ := json.Marshal(tests)
	request, err := http.NewRequest("POST", "", bytes.NewReader(query))
	if err != nil {
		t.Fatal(err)
	}
	result, err := DecodeModelSlice(request, reflect.TypeOf(TestModel{}))
	if err != nil {
		t.Fatal(err)
	}
	checks := result.([]TestModel)
	if checks[0].Key != tests[0].Key || checks[0].Value != tests[0].Value {
		t.Fatal(checks)
	}
}

func TestDecodeModelSlicePtr(t *testing.T) {
	tests := []TestModel{
		{Key: "key", Value: "value"},
	}
	query, _ := json.Marshal(tests)
	request, err := http.NewRequest("POST", "", bytes.NewReader(query))
	if err != nil {
		t.Fatal(err)
	}
	result, err := DecodeModelSlicePtr(request, reflect.TypeOf(TestModel{}))
	if err != nil {
		t.Fatal(err)
	}
	checks := result.(*[]TestModel)
	if (*checks)[0].Key != tests[0].Key || (*checks)[0].Value != tests[0].Value {
		t.Fatal(checks)
	}
}
