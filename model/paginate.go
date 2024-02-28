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

		pageSize, _ := strconv.Atoi(q.Get("limit"))
		switch {
		case pageSize > DefaultPageSize:
			pageSize = DefaultPageSize
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
