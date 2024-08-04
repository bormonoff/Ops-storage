package handlers

import (
	"net/http"
	"strings"

	"ops-storage/internal/core"
)

const TextPlainType string = "text/plain"

func Update(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Only POST methods are allowed", http.StatusMethodNotAllowed)
		return
	}

	contentType, ok := req.Header["Content-Type"]
	if !ok || contentType[0] != TextPlainType {
		http.Error(res, "request should be have a text/plain type", http.StatusBadRequest)
		return
	}

	a := strings.Split(req.URL.Path, "/")
	if len(a) != 5 {
		http.Error(res, "url path has to be in /update/{type}/{name}/{value} format", http.StatusNotFound)
		return
	}
	err := core.Update(a[2], a[3], a[4])
	if err == core.ErrIvalidMetric {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	} else if err != nil {
		http.Error(res, "unexpected error", http.StatusInternalServerError)
		return	
	}
	res.WriteHeader(http.StatusOK)
}
