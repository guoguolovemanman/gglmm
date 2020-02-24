package gglmm

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// MuxParseVarID Mux 解释ID
func MuxParseVarID(r *http.Request) (id int64, err error) {
	return strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
}
