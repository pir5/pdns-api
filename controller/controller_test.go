package controller

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/tj/assert"
)

func dummyContext(t *testing.T, reqType, reqPath string, args interface{}) (echo.Context, *httptest.ResponseRecorder) {
	var rp []byte
	var jsonerr error
	if args != nil {
		switch v := args.(type) {
		case string:
			rp = []byte(v)
		default:
			rp, jsonerr = json.Marshal(v)
			assert.NoError(t, jsonerr)
		}
	}
	req := httptest.NewRequest(reqType, reqPath, bytes.NewReader(rp))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e := echo.New()
	ctx := e.NewContext(req, rec)

	return ctx, rec
}
