package model

import (
	"net/http"
	"strconv"

	"github.com/jinzhu/gorm"
)

const DefaultPageSize = 100

func Paginate(r *http.Request) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		q := r.URL.Query()
		page, _ := strconv.Atoi(q.Get("offset"))
		if page <= 0 {
			page = 1
		}

		pageSize := DefaultPageSize
		requestPageSize, _ := strconv.Atoi(q.Get("limit"))
		if DefaultPageSize > requestPageSize && requestPageSize > 0 {
			pageSize = requestPageSize
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
