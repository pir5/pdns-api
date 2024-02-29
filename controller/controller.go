package controller

import (
	"math"
	"net/http"
	"strconv"

	"github.com/pir5/pdns-api/model"
)

func totalPage(r *http.Request, total int64) int {
	q := r.URL.Query()
	pageSize, _ := strconv.Atoi(q.Get("limit"))
	switch {
	case pageSize > model.DefaultPageSize:
		pageSize = model.DefaultPageSize
	case pageSize <= 0:
		pageSize = 10
	}
	return int(math.Ceil(float64(total) / float64(pageSize)))
}

func setPaginationHeader(w http.ResponseWriter, r *http.Request, total int64) {
	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("offset"))
	if page <= 0 {
		page = 1
	}
	w.Header().Set("X-Pagination-Current-Page", strconv.Itoa(page))
	pageSize, _ := strconv.Atoi(q.Get("limit"))
	switch {
	case pageSize > model.DefaultPageSize:
		pageSize = model.DefaultPageSize
	case pageSize <= 0:
		pageSize = 10
	}
	w.Header().Set("X-Pagination-Limit", strconv.Itoa(pageSize))
	w.Header().Set("X-Pagination-Total-Pages", strconv.Itoa(totalPage(r, total)))
}
