package gglmm

import (
	"encoding/json"
	"net/http"
)

const (
	// ErrEncodeJSON 错误JSON
	ErrEncodeJSON = `{"code":500;"message":"response encode fail"}`
	// ErrUnauthorizedJSON --
	ErrUnauthorizedJSON = `{"code":401,"message":"unauthorized"}`
	// ErrForbiddenJSON --
	ErrForbiddenJSON = `{"code":403,"message":"forbidden"}`
	// SuccessCode 成功码
	SuccessCode = http.StatusOK
	// SuccessMessage 成功消息
	SuccessMessage = "成功"
	// FailCode 失败码
	FailCode = 1000
	// FailMessage 失败消息
	FailMessage = "失败"
)

// Response 响应
type Response struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

// NewResponse 新建响应
func NewResponse() *Response {
	return &Response{}
}

// NewSuccessResponse 新建成功响应
func NewSuccessResponse() *Response {
	return &Response{Code: SuccessCode, Message: SuccessMessage}
}

// NewFailResponse 新建失败响应
func NewFailResponse(failMessage string) *Response {
	return &Response{Code: FailCode, Message: failMessage}
}

// AddData 添加数据
func (response *Response) AddData(key string, value interface{}) *Response {
	if response.Data == nil {
		response.Data = make(map[string]interface{})
	}
	response.Data[key] = value
	return response
}

// WriteJSON 输出JSON
func (response *Response) WriteJSON(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	encoder := json.NewEncoder(w)
	err := encoder.Encode(response)
	if err != nil {
		w.Write([]byte(ErrEncodeJSON))
	}
}

// WriteUnauthorized 输出验证失败
func WriteUnauthorized(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(ErrUnauthorizedJSON))
}

// WriteForbidden 输出无权限
func WriteForbidden(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusForbidden)
	w.Write([]byte(ErrForbiddenJSON))
}

// PageResponse 分页响应
type PageResponse struct {
	List interface{}
	Pagination
}
