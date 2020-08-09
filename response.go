package gglmm

import (
	"encoding/json"
	"net/http"
)

// Resposne --
const (
	ResponseSuccessCode = 0
	ResponseFailCode    = -1
)

// Response 响应
type Response struct {
	StatusCode   int                    `json:"statusCode"`
	ErrorCode    int                    `json:"errorCode"`
	ErrorMessage string                 `json:"errorMessage"`
	Data         map[string]interface{} `json:"data"`
}

// ResponseOf --
func ResponseOf(statusCode int, errorCode int, errorMessage string) *Response {
	return &Response{
		StatusCode:   statusCode,
		ErrorCode:    errorCode,
		ErrorMessage: errorMessage,
	}
}

// UnauthorizedResponse 输出验证失败
func UnauthorizedResponse() *Response {
	return ResponseOf(http.StatusUnauthorized, http.StatusUnauthorized, "")
}

// ForbiddenResponse 输出权限失败
func ForbiddenResponse() *Response {
	return ResponseOf(http.StatusForbidden, http.StatusForbidden, "")
}

// NotFoundResponse 输出不存在
func NotFoundResponse() *Response {
	return ResponseOf(http.StatusNotFound, http.StatusNotFound, "")
}

// OkResponse --
func OkResponse() *Response {
	return ResponseOf(http.StatusOK, ResponseSuccessCode, "")
}

// ErrorResponse --
func ErrorResponse(errorCode int, errorMessage string) *Response {
	return ResponseOf(http.StatusInternalServerError, errorCode, errorMessage)
}

// SucessResponse --
func SucessResponse() *Response {
	return ResponseOf(http.StatusOK, ResponseSuccessCode, "")
}

// FailResponse --
func FailResponse(errorMessage string) *Response {
	return ResponseOf(http.StatusInternalServerError, ResponseFailCode, errorMessage)
}

// AddData 添加数据
func (response *Response) AddData(key string, value interface{}) *Response {
	if response.Data == nil {
		response.Data = make(map[string]interface{})
	}
	response.Data[key] = value
	return response
}

// JSON 输出JSON
func (response Response) JSON(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(response.StatusCode)
	json.NewEncoder(w).Encode(response)
}

// PageResponse 分页响应
type PageResponse struct {
	List interface{}
	Pagination
}
