package gglmm

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

func TestDecodeIDRequest(t *testing.T) {
	url := "/test?id=1&preloads=a,b"
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}
	resultIDRequest := IDRequest{}
	if err = DecodeIDRequest(request, &resultIDRequest); err != nil {
		t.Fatal(err)
	}
	if resultIDRequest.ID != 1 {
		t.Fatal(resultIDRequest)
	}
	if len(resultIDRequest.Preloads) != 2 {
		t.Fatal(resultIDRequest)
	}
	if resultIDRequest.Preloads[0] != "a" || resultIDRequest.Preloads[1] != "b" {
		t.Fatal(resultIDRequest)
	}
}

func TestDecodeFilterRequest(t *testing.T) {
	filterRequest := FilterRequest{
		Filters: []*Filter{
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
	resultfilterRequest := FilterRequest{}
	if err = DecodeBody(request, &resultfilterRequest); err != nil {
		t.Fatal(err)
	}
	if resultfilterRequest.Order != "id" {
		t.Fatal(resultfilterRequest)
	}
	filter := resultfilterRequest.Filters[0]
	if filter.Field != "A" || filter.Operate != FilterOperateEqual || filter.Value != "B" {
		t.Fatal(resultfilterRequest)
	}
}

func TestDecodePageRequest(t *testing.T) {
	pageRequest := PageRequest{
		FilterRequest: FilterRequest{
			Filters: []*Filter{
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
	resultPageRequest := PageRequest{}
	if err = DecodeBody(request, &resultPageRequest); err != nil {
		t.Fatal(err)
	}
	filterRequest := resultPageRequest.FilterRequest
	if filterRequest.Order != "id" {
		t.Fatal(resultPageRequest)
	}
	filter := filterRequest.Filters[0]
	if filter.Field != "A" || filter.Operate != FilterOperateEqual || filter.Value != "B" {
		t.Fatal(resultPageRequest)
	}
	pagination := resultPageRequest.Pagination
	if pagination.PageSize != 0 || pagination.PageIndex != 1 || pagination.Total != 2 {
		t.Fatal(resultPageRequest)
	}
}
