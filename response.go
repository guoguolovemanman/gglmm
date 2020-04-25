package gglmm

import (
	"encoding/json"
	"net/http"
)

const (
	// ResposneCodeFail --
	ResposneCodeFail = 1000
	// ResponseCodeRequestError --
	ResponseCodeRequestError = 1001
)

// Response 响应
type Response struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

// ResponseOf --
func ResponseOf(code int, message string) *Response {
	return &Response{
		Code:    code,
		Message: message,
	}
}

// UnauthorizedResponse 输出验证失败
func UnauthorizedResponse() *Response {
	return ResponseOf(http.StatusUnauthorized, "")
}

// ForbiddenResponse 输出权限失败
func ForbiddenResponse() *Response {
	return ResponseOf(http.StatusForbidden, "")
}

// NotFoundResponse 输出不存在
func NotFoundResponse() *Response {
	return ResponseOf(http.StatusNotFound, "")
}

// OkResponse --
func OkResponse() *Response {
	return ResponseOf(http.StatusOK, "")
}

// InternalErrorResponse --
func InternalErrorResponse(message string) *Response {
	return ResponseOf(http.StatusInternalServerError, message)
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
func (response *Response) JSON(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if response.Code < ResposneCodeFail {
		w.WriteHeader(response.Code)
	}
	json.NewEncoder(w).Encode(response)
}

// PageResponse 分页响应
type PageResponse struct {
	List interface{}
	Pagination
}
