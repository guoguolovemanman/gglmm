package gglmm

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
)

// DecodeIDRequest --
func DecodeIDRequest(r *http.Request, idRequest *IDRequest) error {
	id, err := PathVarID(r)
	if err != nil {
		idQuery := r.FormValue("id")
		if idQuery == "" {
			return ErrRequest
		}
		id, err = strconv.ParseUint(idQuery, 10, 64)
		if err != nil {
			return ErrRequest
		}
	}
	idRequest.ID = id
	preloadsQuery, err := PathVar(r, "preloads")
	if err != nil {
		preloadsQuery = r.FormValue("preloads")
	}
	if preloadsQuery != "" {
		idRequest.Preloads = strings.Split(preloadsQuery, ",")
	}
	return nil
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
