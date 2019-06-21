package pdns_api

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"

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
	}

	return ds, nil
}
func (d *domainModelStub) UpdateByName(name string, newDomain *model.Domain) (bool, error) {
	// notfound			return false, nil
	// err		return false, r.Error
	return true, nil
}

func (d *domainModelStub) DeleteByName(name string) (bool, error) {
	return true, nil
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
