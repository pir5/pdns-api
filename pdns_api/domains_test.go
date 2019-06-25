package pdns_api

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/pir5/pdns-api/model"
)

type domainModelStub struct {
}

func (d *domainModelStub) FindBy(params map[string]interface{}) (model.Domains, error) {
	ds := model.Domains{}

	switch params["name"].([]string)[0] {
	case "ok":
		ds = model.Domains{
			model.Domain{
				ID:   1,
				Name: "test.com",
			},
		}
	case "error_please":
		return nil, fmt.Errorf("give error to you")
	case "allowed":
		if params["name"].([]string)[1] == "allow.com" {
			ds = model.Domains{
				model.Domain{
					ID:   1,
					Name: "allow.com",
				},
			}
		} else {
			return nil, fmt.Errorf("give error to you")
		}
	}

	return ds, nil
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

func (d *domainModelStub) Create(newDomain *model.Domain) error {
	return nil
}

func Test_domainHandler_getDomains(t *testing.T) {
	type fields struct {
		domainModel model.DomainModel
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
			wantCode:  http.StatusNotFound,
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
		{
			name: "allowed domain",
			fields: fields{
				domainModel: &domainModelStub{},
			},
			wantErr:   false,
			wantCode:  http.StatusOK,
			queryName: "allowed",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &domainHandler{
				domainModel: tt.fields.domainModel,
			}

			globalConfig = Config{
				TokenAuth: tokenAuth{
					AuthType: AuthTypeHTTP,
				},
			}

			q := make(url.Values)
			q.Set("name", tt.queryName)
			ctx, rec := dummyContext(t, "GET", "/domains?"+q.Encode(), nil)

			ctx.Set(AllowDomainsKey, []string{"allow.com"})
			if err := h.getDomains(ctx); (err != nil) != tt.wantErr {
				t.Errorf("domainHandler.getDomains() error = %v, wantErr %v", err, tt.wantErr)
			}
			if rec.Code != tt.wantCode {
				t.Errorf("domainHandler.getDomains() got different http status code = %d, wantCode %d", rec.Code, tt.wantCode)
			}
		})
	}
}

func Test_domainHandler_updateDomain(t *testing.T) {
	type fields struct {
		domainModel model.DomainModel
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
				Name: "ok.com",
				ID:   9999,
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
				Name: "ok.com",
				ID:   1111,
			},
			wantErr:   false,
			wantCode:  http.StatusNotFound,
			queryName: "notfound.com",
		},
		{
			name: "not allowed",
			fields: fields{
				domainModel: &domainModelStub{},
			},
			wantErr:   false,
			wantCode:  http.StatusForbidden,
			queryName: "notallowed.com",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &domainHandler{
				domainModel: tt.fields.domainModel,
			}

			ctx, rec := dummyContext(t, "PUT", "/domains/:", tt.args)
			ctx.SetParamNames("name")
			ctx.SetParamValues(tt.queryName)

			ctx.Set(AllowDomainsKey, []string{"ok.com", "notfound.com"})
			if err := h.updateDomain(ctx); (err != nil) != tt.wantErr {
				t.Errorf("domainHandler.updateDomain() error = %v, wantErr %v", err, tt.wantErr)
			}

			if rec.Code != tt.wantCode {
				t.Errorf("domainHandler.updateDomains() got different http status code = %d, wantCode %d", rec.Code, tt.wantCode)
			}
		})
	}
}

func Test_domainHandler_deleteDomain(t *testing.T) {
	type fields struct {
		domainModel model.DomainModel
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
		{
			name: "not allowed",
			fields: fields{
				domainModel: &domainModelStub{},
			},
			wantErr:   false,
			wantCode:  http.StatusForbidden,
			queryName: "notallowed.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &domainHandler{
				domainModel: tt.fields.domainModel,
			}

			ctx, rec := dummyContext(t, "DELETE", "/domains/:", nil)
			ctx.SetParamNames("name")
			ctx.SetParamValues(tt.queryName)

			ctx.Set(AllowDomainsKey, []string{"ok.com", "notfound.com"})
			if err := h.deleteDomain(ctx); (err != nil) != tt.wantErr {
				t.Errorf("domainHandler.deleteDomain() error = %v, wantErr %v", err, tt.wantErr)
			}

			if rec.Code != tt.wantCode {
				t.Errorf("domainHandler.deleteDomains() got different http status code = %d, wantCode %d", rec.Code, tt.wantCode)
			}

		})
	}
}

func Test_domainHandler_createDomain(t *testing.T) {
	type fields struct {
		domainModel model.DomainModel
	}
	type args struct {
		c echo.Context
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
				Name: "ok.com",
				ID:   9999,
			},
			wantErr:   false,
			wantCode:  http.StatusCreated,
			queryName: "ok.com",
		},
		{
			name: "not allowed",
			fields: fields{
				domainModel: &domainModelStub{},
			},
			wantErr:   false,
			wantCode:  http.StatusForbidden,
			queryName: "notallowed.com",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &domainHandler{
				domainModel: tt.fields.domainModel,
			}

			ctx, rec := dummyContext(t, "PUT", "/domains/:", tt.args)
			ctx.SetParamNames("name")
			ctx.SetParamValues(tt.queryName)

			ctx.Set(AllowDomainsKey, []string{"ok.com"})
			if err := h.createDomain(ctx); (err != nil) != tt.wantErr {
				t.Errorf("domainHandler.createDomain() error = %v, wantErr %v", err, tt.wantErr)
			}

			if rec.Code != tt.wantCode {
				t.Errorf("domainHandler.createDomains() got different http status code = %d, wantCode %d", rec.Code, tt.wantCode)
			}
		})
	}

}