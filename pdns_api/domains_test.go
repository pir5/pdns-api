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

func init() {
	globalConfig = Config{
		Auth: auth{
			AuthType: AuthTypeHTTP,
		},
	}
}

func (d *domainModelStub) FindBy(params map[string]interface{}) (model.Domains, error) {
	ds := model.Domains{}

	if params["name"] != nil {
		switch params["name"].([]string)[0] {
		case "ok":
			ds = model.Domains{
				model.Domain{
					ID:   1,
					Name: "ok.com",
				},
			}
		case "error_please":
			return nil, fmt.Errorf("give error to you")
		case "deny":
			ds = model.Domains{
				model.Domain{
					ID:   1,
					Name: "deny.com",
				},
			}
		}
	} else if params["id"] != nil {
		switch params["id"].(int) {
		case 1, 4:
			ds = model.Domains{
				model.Domain{
					ID:   1,
					Name: "ok.com",
				},
			}
		case 2:
			return nil, fmt.Errorf("give error to you")
		case 3:
			ds = model.Domains{
				model.Domain{
					ID:   1,
					Name: "deny.com",
				},
			}
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
		{
			name: "deny",
			fields: fields{
				domainModel: &domainModelStub{},
			},
			wantErr:   false,
			wantCode:  http.StatusOK,
			queryName: "deny",
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

			ctx.Set(AllowDomainsKey, []string{"ok.com"})
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
				Name: "ok.com",
				ID:   9999,
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
				Name: "ok.com",
				ID:   1111,
			},
			wantErr:   false,
			wantCode:  http.StatusNotFound,
			queryName: "4",
		},
		{
			name: "deny",
			fields: fields{
				domainModel: &domainModelStub{},
			},
			wantErr:   false,
			wantCode:  http.StatusForbidden,
			queryName: "3",
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

			ctx.Set(AllowDomainsKey, []string{"ok.com", "notfound.com"})
			if err := h.updateDomainByID(ctx); (err != nil) != tt.wantErr {
				t.Errorf("domainHandler.updateDomainByID() error = %v, wantErr %v", err, tt.wantErr)
			}

			if rec.Code != tt.wantCode {
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
			name: "deny",
			fields: fields{
				domainModel: &domainModelStub{},
			},
			wantErr:   false,
			wantCode:  http.StatusForbidden,
			queryName: "deny",
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

			ctx.Set(AllowDomainsKey, []string{"ok.com", "notfound.com"})
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
		{
			name: "deny",
			fields: fields{
				domainModel: &domainModelStub{},
			},
			wantErr:   false,
			wantCode:  http.StatusForbidden,
			queryName: "deny",
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

			ctx.Set(AllowDomainsKey, []string{"ok.com", "notfound.com"})
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
		{
			name: "deny",
			fields: fields{
				domainModel: &domainModelStub{},
			},
			wantErr:  false,
			wantCode: http.StatusForbidden,
			queryID:  "3",
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

			ctx.Set(AllowDomainsKey, []string{"ok.com", "notfound.com"})
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
			name: "deny",
			fields: fields{
				domainModel: &domainModelStub{},
			},
			args: model.Domain{
				Name: "deny.com",
				ID:   1,
			},
			wantErr:   false,
			wantCode:  http.StatusForbidden,
			queryName: "deny.com",
		},
		{
			name: "invalid domain name",
			fields: fields{
				domainModel: &domainModelStub{},
			},
			args: model.Domain{
				Name: "@.com",
				ID:   1,
			},
			wantErr:   false,
			wantCode:  http.StatusBadRequest,
			queryName: "@.com",
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
