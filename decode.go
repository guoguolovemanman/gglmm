package gglmm

import (
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

// DecodeIDRequest 解码ID请求
func DecodeIDRequest(r *http.Request) (*IDRequest, error) {
	id, err := PathVarID(r)
	if err != nil {
		idQuery := r.FormValue("id")
		if idQuery == "" {
			return nil, ErrRequest
		}
		id, err = strconv.ParseInt(idQuery, 10, 64)
		if err != nil {
			return nil, ErrRequest
		}
	}
	idRequest := &IDRequest{
		ID: id,
	}
	preloadsQuery := r.FormValue("preloads")
	if preloadsQuery != "" {
		idRequest.Preloads = strings.Split(preloadsQuery, ",")
	}
	return idRequest, nil
}

// DecodeBody 解码请求体
func DecodeBody(r *http.Request, body interface{}) error {
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(body); err != nil {
		if err != io.EOF {
			return ErrRequest
		}
	}
	return nil
}

// DecodeFilterRequest 解码过滤请求
func DecodeFilterRequest(r *http.Request) (*FilterRequest, error) {
	filterRequest := &FilterRequest{}
	err := DecodeBody(r, filterRequest)
	return filterRequest, err
}

// DecodePageRequest 解码分页请求
func DecodePageRequest(r *http.Request) (*PageRequest, error) {
	pageRequest := &PageRequest{
		Pagination: Pagination{
			PageSize:  DefaultPageSize,
			PageIndex: FirstPageIndex,
		},
	}
	err := DecodeBody(r, pageRequest)
	return pageRequest, err
}

// DecodeModelPtr 解码模型指针
func DecodeModelPtr(r *http.Request, modelType reflect.Type) (interface{}, error) {
	decoder := json.NewDecoder(r.Body)
	model := reflect.New(modelType).Interface()
	if err := decoder.Decode(model); err != nil {
		return nil, err
	}
	return model, nil
}

// DecodeModel 解码模型
func DecodeModel(r *http.Request, modelType reflect.Type) (interface{}, error) {
	decoder := json.NewDecoder(r.Body)
	model := reflect.New(modelType)
	if err := decoder.Decode(model.Interface()); err != nil {
		return nil, err
	}
	return model.Elem().Interface(), nil
}

// DecodeModelSlicePtr 解码模型列表指针
func DecodeModelSlicePtr(r *http.Request, modelType reflect.Type) (interface{}, error) {
	decoder := json.NewDecoder(r.Body)
	list := reflect.New(reflect.SliceOf(modelType))
	if err := decoder.Decode(list.Interface()); err != nil {
		return nil, err
	}
	return list, nil
}

// DecodeModelSlice 解码模型列表
func DecodeModelSlice(r *http.Request, modelType reflect.Type) (interface{}, error) {
	decoder := json.NewDecoder(r.Body)
	list := reflect.New(reflect.SliceOf(modelType))
	if err := decoder.Decode(list.Interface()); err != nil {
		return nil, err
	}
	return list.Elem().Interface(), nil
}
