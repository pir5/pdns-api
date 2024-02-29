package controller

import (
	"database/sql"
	"database/sql/driver"
	"math"
	"net/http"
	"reflect"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/pir5/pdns-api/model"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
	validate.RegisterCustomTypeFunc(ValidateValuer, sql.NullString{}, sql.NullInt64{}, sql.NullInt32{})
}

func ValidateValuer(field reflect.Value) interface{} {

	if valuer, ok := field.Interface().(driver.Valuer); ok {

		val, err := valuer.Value()
		if err == nil {
			return val
		}
		// handle the error how you want
	}

	return nil
}

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
