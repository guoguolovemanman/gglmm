package gglmm

import (
	"encoding/json"
	"io"
	"net/http"
	"reflect"
)

// Decode 解码请求体
func Decode(r *http.Request, body interface{}) error {
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(body); err != nil {
		return err
	}
	return nil
}

// DecodeFilterRequest 解码过滤请求
func DecodeFilterRequest(r *http.Request) (FilterRequest, error) {
	decoder := json.NewDecoder(r.Body)
	filterRequest := FilterRequest{}
	if err := decoder.Decode(&filterRequest); err != nil {
		if err != io.EOF {
			return filterRequest, err
		}
	}
	return filterRequest, nil
}

// DecodePageRequest 解码分页请求
func DecodePageRequest(r *http.Request) (PageRequest, error) {
	decoder := json.NewDecoder(r.Body)
	pageRequest := PageRequest{
		Pagination: Pagination{
			PageSize:  DefaultPageSize,
			PageIndex: FirstPageIndex,
		},
	}
	if err := decoder.Decode(&pageRequest); err != nil {
		if err != io.EOF {
			return pageRequest, err
		}
	}
	return pageRequest, nil
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
	list := reflect.New(reflect.SliceOf(modelType)).Interface()
	if err := decoder.Decode(list); err != nil {
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
