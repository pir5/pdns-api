package controller

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/pir5/pdns-api/model"
	null "gopkg.in/guregu/null.v3"
)

type domainModelStub struct {
}

func (d *domainModelStub) FindBy(req *http.Request, params map[string]interface{}) (model.Domains, int64, error) {
	ds := model.Domains{}
	totalPage := int64(1)
	if params["name"] != nil {
		switch params["name"].([]string)[0] {
		case "ok":
			ds = model.Domains{
				model.Domain{
					ID:   1,
					Name: null.NewString("ok.com", true),
				},
			}
		case "error_please":
			return nil, 0, fmt.Errorf("give error to you")
		}
	} else if params["id"] != nil {
		switch params["id"].(int) {
		case 1, 4:
			ds = model.Domains{
				model.Domain{
					ID:   1,
					Name: null.NewString("ok.com", true),
				},
			}
		case 2:
			return nil, 0, fmt.Errorf("give error to you")
		case 3:
			ds = model.Domains{
				model.Domain{
					ID:   1,
					Name: null.NewString("deny.com", true),
				},
			}
		}
	}

	return ds, totalPage, nil
}
func (d *domainModelStub) UpdateByName(name string, newDomain *model.Domain) (bool, error) {
	switch name {
	case "ok.com":
		return true, nil
	}
	return false, nil
}

func (d *domainModelStub) DeleteByName(name string) (bool, error) {
	switch name {
	case "ok.com":
		return true, nil
	}
	return false, nil
}

func (d *domainModelStub) UpdateByID(id string, newDomain *model.Domain) (bool, error) {
	switch id {
	case "1":
		return true, nil
	case "4":
		return false, nil
	}
	return false, nil
}

func (d *domainModelStub) DeleteByID(id string) (bool, error) {
	switch id {
	case "1":
		return true, nil
	case "2":
		return false, nil
	}
	return false, nil
}

func (d *domainModelStub) Create(newDomain *model.Domain) error {
	return nil
}

func Test_domainHandler_getDomains(t *testing.T) {
	type fields struct {
		domainModel model.DomainModeler
	}
	tests := []struct {
		name      string
		fields    fields
		wantErr   bool
		wantCode  int
		queryName string
	}{
		{
			name: "ok",
			fields: fields{
				domainModel: &domainModelStub{},
			},
			wantErr:   false,
			wantCode:  http.StatusOK,
			queryName: "ok",
		},
		{
			name: "notfound",
			fields: fields{
				domainModel: &domainModelStub{},
			},
			wantErr:   false,
			wantCode:  http.StatusOK,
			queryName: "notfound",
		},
		{
			name: "get error",
			fields: fields{
				domainModel: &domainModelStub{},
			},
			wantErr:   false,
			wantCode:  http.StatusInternalServerError,
			queryName: "error_please",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &domainHandler{
				domainModel: tt.fields.domainModel,
			}

			q := make(url.Values)
			q.Set("name", tt.queryName)
			ctx, rec := dummyContext(t, "GET", "/domains?"+q.Encode(), nil)

			if err := h.getDomains(ctx); (err != nil) != tt.wantErr {
				t.Errorf("domainHandler.getDomains() error = %v, wantErr %v", err, tt.wantErr)
			}
			if rec.Code != tt.wantCode {
				t.Errorf("domainHandler.getDomains() got different http status code = %d, wantCode %d", rec.Code, tt.wantCode)
			}
		})
	}
}

func Test_domainHandler_updateDomainByID(t *testing.T) {
	type fields struct {
		domainModel model.DomainModeler
	}
	tests := []struct {
		name      string
		fields    fields
		args      model.Domain
		wantErr   bool
		wantCode  int
		queryName string
	}{
		{
			name: "ok",
			fields: fields{
				domainModel: &domainModelStub{},
			},
			args: model.Domain{
				Name: null.NewString("ok.com", true),
				ID:   9999,
				Type: null.NewString("NATIVE", true),
			},
			wantErr:   false,
			wantCode:  http.StatusOK,
			queryName: "1",
		},
		{
			name: "not found",
			fields: fields{
				domainModel: &domainModelStub{},
			},
			args: model.Domain{
				Name: null.NewString("ok.com", true),
				ID:   1111,
				Type: null.NewString("NATIVE", true),
			},
			wantErr:   false,
			wantCode:  http.StatusNotFound,
			queryName: "4",
		},
		{
			name: "invalid domain name",
			fields: fields{
				domainModel: &domainModelStub{},
			},
			args: model.Domain{
				Name: null.NewString("@.com", true),
				ID:   1,
				Type: null.NewString("NATIVE", true),
			},
			wantErr:   false,
			wantCode:  http.StatusBadRequest,
			queryName: "1",
		},
		{
			name: "invalid domain type",
			fields: fields{
				domainModel: &domainModelStub{},
			},
			args: model.Domain{
				Name: null.NewString("ok.com", true),
				ID:   1,
				Type: null.NewString("NG", true),
			},
			wantErr:   false,
			wantCode:  http.StatusBadRequest,
			queryName: "1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &domainHandler{
				domainModel: tt.fields.domainModel,
			}

			ctx, rec := dummyContext(t, "PUT", "/domains/:", tt.args)
			ctx.SetParamNames("id")
			ctx.SetParamValues(tt.queryName)

			if err := h.updateDomainByID(ctx); (err != nil) != tt.wantErr {
				t.Errorf("domainHandler.updateDomainByID() error = %v, wantErr %v", err, tt.wantErr)
			}

			if rec.Code != tt.wantCode {
				t.Errorf("%+v", rec)
				t.Errorf("domainHandler.updateDomainsByID() got different http status code = %d, wantCode %d", rec.Code, tt.wantCode)
			}
		})
	}
}

func Test_domainHandler_updateDomainByName(t *testing.T) {
	type fields struct {
		domainModel model.DomainModeler
	}
	tests := []struct {
		name      string
		fields    fields
		args      model.Domain
		wantErr   bool
		wantCode  int
		queryName string
	}{
		{
			name: "ok",
			fields: fields{
				domainModel: &domainModelStub{},
			},
			args: model.Domain{
				Name: null.NewString("ok.com", true),
				ID:   9999,
				Type: null.NewString("NATIVE", true),
			},
			wantErr:   false,
			wantCode:  http.StatusOK,
			queryName: "ok.com",
		},
		{
			name: "not found",
			fields: fields{
				domainModel: &domainModelStub{},
			},
			args: model.Domain{
				Name: null.NewString("ok.com", true),
				ID:   1111,
				Type: null.NewString("NATIVE", true),
			},
			wantErr:   false,
			wantCode:  http.StatusNotFound,
			queryName: "notfound.com",
		},
		{
			name: "invalid domain name",
			fields: fields{
				domainModel: &domainModelStub{},
			},
			args: model.Domain{
				Name: null.NewString("@.com", true),
				ID:   1,
				Type: null.NewString("NATIVE", true),
			},
			wantErr:   false,
			wantCode:  http.StatusBadRequest,
			queryName: "ok.com",
		},
		{
			name: "invalid domain type",
			fields: fields{
				domainModel: &domainModelStub{},
			},
			args: model.Domain{
				Name: null.NewString("ok.com", true),
				ID:   1,
				Type: null.NewString("NG", true),
			},
			wantErr:   false,
			wantCode:  http.StatusBadRequest,
			queryName: "ok.com",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &domainHandler{
				domainModel: tt.fields.domainModel,
			}

			ctx, rec := dummyContext(t, "PUT", "/domains/name/:", tt.args)
			ctx.SetParamNames("name")
			ctx.SetParamValues(tt.queryName)

			if err := h.updateDomainByName(ctx); (err != nil) != tt.wantErr {
				t.Errorf("domainHandler.updateDomainByName() error = %v, wantErr %v", err, tt.wantErr)
			}

			if rec.Code != tt.wantCode {
				t.Errorf("domainHandler.updateDomainsByName() got different http status code = %d, wantCode %d", rec.Code, tt.wantCode)
			}
		})
	}
}

func Test_domainHandler_deleteDomainByName(t *testing.T) {
	type fields struct {
		domainModel model.DomainModeler
	}
	tests := []struct {
		name      string
		fields    fields
		wantErr   bool
		wantCode  int
		queryName string
	}{
		{
			name: "ok",
			fields: fields{
				domainModel: &domainModelStub{},
			},
			wantErr:   false,
			wantCode:  http.StatusNoContent,
			queryName: "ok.com",
		},
		{
			name: "not found",
			fields: fields{
				domainModel: &domainModelStub{},
			},
			wantErr:   false,
			wantCode:  http.StatusNotFound,
			queryName: "notfound.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &domainHandler{
				domainModel: tt.fields.domainModel,
			}

			ctx, rec := dummyContext(t, "DELETE", "/domains/name/:", nil)
			ctx.SetParamNames("name")
			ctx.SetParamValues(tt.queryName)

			if err := h.deleteDomainByName(ctx); (err != nil) != tt.wantErr {
				t.Errorf("domainHandler.deleteDomainByName() error = %v, wantErr %v", err, tt.wantErr)
			}

			if rec.Code != tt.wantCode {
				t.Errorf("domainHandler.deleteDomainsByName() got different http status code = %d, wantCode %d", rec.Code, tt.wantCode)
			}

		})
	}
}

func Test_domainHandler_deleteDomainByID(t *testing.T) {
	type fields struct {
		domainModel model.DomainModeler
	}
	tests := []struct {
		name     string
		fields   fields
		wantErr  bool
		wantCode int
		queryID  string
	}{
		{
			name: "ok",
			fields: fields{
				domainModel: &domainModelStub{},
			},
			wantErr:  false,
			wantCode: http.StatusNoContent,
			queryID:  "1",
		},
		{
			name: "not found",
			fields: fields{
				domainModel: &domainModelStub{},
			},
			wantErr:  false,
			wantCode: http.StatusNotFound,
			queryID:  "4",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &domainHandler{
				domainModel: tt.fields.domainModel,
			}

			ctx, rec := dummyContext(t, "DELETE", "/domains/id/:", nil)
			ctx.SetParamNames("id")
			ctx.SetParamValues(tt.queryID)

			if err := h.deleteDomainByID(ctx); (err != nil) != tt.wantErr {
				t.Errorf("domainHandler.deleteDomainByID() error = %v, wantErr %v", err, tt.wantErr)
			}

			if rec.Code != tt.wantCode {
				t.Errorf("domainHandler.deleteDomainsByID() got different http status code = %d, wantCode %d", rec.Code, tt.wantCode)
			}

		})
	}
}

func Test_domainHandler_createDomain(t *testing.T) {
	type fields struct {
		domainModel model.DomainModeler
	}
	tests := []struct {
		name      string
		fields    fields
		args      model.Domain
		wantErr   bool
		wantCode  int
		queryName string
	}{
		{
			name: "ok",
			fields: fields{
				domainModel: &domainModelStub{},
			},
			args: model.Domain{
				Name: null.NewString("ok.com", true),
				ID:   9999,
				Type: null.NewString("NATIVE", true),
			},
			wantErr:   false,
			wantCode:  http.StatusCreated,
			queryName: "ok.com",
		},
		{
			name: "invalid domain name",
			fields: fields{
				domainModel: &domainModelStub{},
			},
			args: model.Domain{
				Name: null.NewString("@.com", true),
				ID:   1,
				Type: null.NewString("NATIVE", true),
			},
			wantErr:   false,
			wantCode:  http.StatusBadRequest,
			queryName: "@.com",
		},
		{
			name: "invalid domain type",
			fields: fields{
				domainModel: &domainModelStub{},
			},
			args: model.Domain{
				Name: null.NewString("example.com", true),
				ID:   1,
				Type: null.NewString("NG", true),
			},
			wantErr:   false,
			wantCode:  http.StatusBadRequest,
			queryName: "example.com",
		},
		{
			name: "empty domain type",
			fields: fields{
				domainModel: &domainModelStub{},
			},
			args: model.Domain{
				Name: null.NewString("example.com", true),
				ID:   1,
			},
			wantErr:   false,
			wantCode:  http.StatusBadRequest,
			queryName: "example.com",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &domainHandler{
				domainModel: tt.fields.domainModel,
			}

			ctx, rec := dummyContext(t, "POST", "/domains/:", tt.args)
			ctx.SetParamNames("name")
			ctx.SetParamValues(tt.queryName)

			if err := h.createDomain(ctx); (err != nil) != tt.wantErr {
				t.Errorf("domainHandler.createDomain() error = %v, wantErr %v", err, tt.wantErr)
			}

			if rec.Code != tt.wantCode {
				t.Errorf("domainHandler.createDomains() got different http status code = %d, wantCode %d", rec.Code, tt.wantCode)
			}
		})
	}

}
