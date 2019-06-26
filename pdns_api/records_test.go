package pdns_api

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/pir5/pdns-api/model"
)

type recordModelStub struct {
}

func (d *recordModelStub) FindBy(params map[string]interface{}) (model.Records, error) {
	ds := model.Records{}

	switch params["id"].([]string)[0] {
	case "1":
		ds = model.Records{
			model.Record{
				ID:       1,
				Name:     "ok.com",
				DomainID: 1,
				Domain: model.Domain{
					Name: "ok.com",
				},
			},
		}
	case "2":
		return nil, fmt.Errorf("give error to you")
	case "3":
		ds = model.Records{
			model.Record{
				ID:       1,
				Name:     "deny.com",
				DomainID: 3,
				Domain: model.Domain{
					Name: "deny.com",
				},
			},
		}
	}

	return ds, nil
}
func (d *recordModelStub) UpdateByID(id string, newRecord *model.Record) (bool, error) {
	switch id {
	case "1":
		return true, nil
	}
	return false, nil
}

func (d *recordModelStub) DeleteByID(id string) (bool, error) {
	switch id {
	case "1":
		return true, nil
	}
	return false, nil
}

func (d *recordModelStub) Create(newRecord *model.Record) error {
	return nil
}
func Test_recordHandler_getRecords(t *testing.T) {
	type fields struct {
		recordModel model.RecordModel
		domainModel model.DomainModel
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
				recordModel: &recordModelStub{},
			},
			wantErr:  false,
			wantCode: http.StatusOK,
			queryID:  "1",
		},
		{
			name: "notfound",
			fields: fields{
				domainModel: &domainModelStub{},
				recordModel: &recordModelStub{},
			},
			wantErr:  false,
			wantCode: http.StatusNotFound,
			queryID:  "0",
		},
		{
			name: "get error",
			fields: fields{
				domainModel: &domainModelStub{},
				recordModel: &recordModelStub{},
			},
			wantErr:  false,
			wantCode: http.StatusInternalServerError,
			queryID:  "2",
		},
		{
			name: "deny",
			fields: fields{
				domainModel: &domainModelStub{},
				recordModel: &recordModelStub{},
			},
			wantErr:  false,
			wantCode: http.StatusNotFound,
			queryID:  "3",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &recordHandler{
				recordModel: tt.fields.recordModel,
				domainModel: tt.fields.domainModel,
			}

			q := make(url.Values)
			q.Set("id", tt.queryID)
			ctx, rec := dummyContext(t, "GET", "/records?"+q.Encode(), nil)
			ctx.Set(AllowDomainsKey, []string{"ok.com", "allow.com"})

			if err := h.getRecords(ctx); (err != nil) != tt.wantErr {
				t.Errorf("recordHandler.getRecords() error = %v, wantErr %v", err, tt.wantErr)
			}
			if rec.Code != tt.wantCode {
				t.Errorf("recordHandler.getRecords() got different http status code = %d, wantCode %d", rec.Code, tt.wantCode)
			}
		})
	}
}

func Test_recordHandler_updateRecord(t *testing.T) {
	type fields struct {
		recordModel model.RecordModel
		domainModel model.DomainModel
	}
	tests := []struct {
		name     string
		fields   fields
		args     model.Record
		wantErr  bool
		wantCode int
		queryID  string
	}{
		{
			name: "ok",
			fields: fields{
				domainModel: &domainModelStub{},
				recordModel: &recordModelStub{},
			},
			args: model.Record{
				Name: "ok.com",
				ID:   9999,
			},
			wantErr:  false,
			wantCode: http.StatusOK,
			queryID:  "1",
		},
		{
			name: "not found",
			fields: fields{
				domainModel: &domainModelStub{},
				recordModel: &recordModelStub{},
			},
			args: model.Record{
				Name: "ok.com",
				ID:   1111,
			},
			wantErr:  false,
			wantCode: http.StatusNotFound,
			queryID:  "0",
		},
		{
			name: "deny",
			fields: fields{
				domainModel: &domainModelStub{},
				recordModel: &recordModelStub{},
			},
			wantErr:  false,
			wantCode: http.StatusForbidden,
			queryID:  "3",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &recordHandler{
				recordModel: tt.fields.recordModel,
				domainModel: tt.fields.domainModel,
			}
			ctx, rec := dummyContext(t, "PUT", "/records/:", tt.args)
			ctx.SetParamNames("id")
			ctx.SetParamValues(tt.queryID)

			ctx.Set(AllowDomainsKey, []string{"ok.com", "notfound.com"})
			if err := h.updateRecord(ctx); (err != nil) != tt.wantErr {
				t.Errorf("recordHandler.updateRecord() error = %v, wantErr %v", err, tt.wantErr)
			}

			if rec.Code != tt.wantCode {
				t.Errorf("recordHandler.updateRecords() got different http status code = %d, wantCode %d", rec.Code, tt.wantCode)
			}

		})
	}
}

func Test_recordHandler_deleteRecord(t *testing.T) {
	type fields struct {
		recordModel model.RecordModel
		domainModel model.DomainModel
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
				recordModel: &recordModelStub{},
			},
			wantErr:  false,
			wantCode: http.StatusNoContent,
			queryID:  "1",
		},
		{
			name: "not found",
			fields: fields{
				domainModel: &domainModelStub{},
				recordModel: &recordModelStub{},
			},
			wantErr:  false,
			wantCode: http.StatusNotFound,
			queryID:  "0",
		},
		{
			name: "deny",
			fields: fields{
				domainModel: &domainModelStub{},
				recordModel: &recordModelStub{},
			},
			wantErr:  false,
			wantCode: http.StatusForbidden,
			queryID:  "3",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &recordHandler{
				recordModel: tt.fields.recordModel,
				domainModel: tt.fields.domainModel,
			}
			ctx, rec := dummyContext(t, "DELETE", "/records/:", nil)
			ctx.SetParamNames("id")
			ctx.SetParamValues(tt.queryID)

			ctx.Set(AllowDomainsKey, []string{"ok.com", "notfound.com"})
			if err := h.deleteRecord(ctx); (err != nil) != tt.wantErr {
				t.Errorf("recordHandler.deleteRecord() error = %v, wantErr %v", err, tt.wantErr)
			}
			if rec.Code != tt.wantCode {
				t.Errorf("recordHandler.deleteRecords() got different http status code = %d, wantCode %d", rec.Code, tt.wantCode)
			}

		})
	}
}

func Test_recordHandler_createRecord(t *testing.T) {
	type fields struct {
		recordModel model.RecordModel
		domainModel model.DomainModel
	}
	tests := []struct {
		name     string
		fields   fields
		args     model.Record
		wantErr  bool
		wantCode int
	}{
		{
			name: "ok",
			fields: fields{
				domainModel: &domainModelStub{},
				recordModel: &recordModelStub{},
			},
			args: model.Record{
				Name:     "ok.com",
				ID:       9999,
				DomainID: 1,
			},
			wantErr:  false,
			wantCode: http.StatusCreated,
		},
		{
			name: "not found domain",
			fields: fields{
				domainModel: &domainModelStub{},
				recordModel: &recordModelStub{},
			},
			args: model.Record{
				Name:     "ok.com",
				ID:       1111,
				DomainID: 0,
			},
			wantErr:  false,
			wantCode: http.StatusNotFound,
		},
		{
			name: "deny",
			fields: fields{
				domainModel: &domainModelStub{},
				recordModel: &recordModelStub{},
			},
			args: model.Record{
				Name:     "ok.com",
				ID:       9999,
				DomainID: 3,
			},
			wantErr:  false,
			wantCode: http.StatusForbidden,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &recordHandler{
				recordModel: tt.fields.recordModel,
				domainModel: tt.fields.domainModel,
			}
			ctx, rec := dummyContext(t, "POST", "/records", tt.args)

			ctx.Set(AllowDomainsKey, []string{"ok.com", "notfound.com"})
			if err := h.createRecord(ctx); (err != nil) != tt.wantErr {
				t.Errorf("recordHandler.createRecord() error = %v, wantErr %v", err, tt.wantErr)
			}

			if rec.Code != tt.wantCode {
				t.Errorf("recordHandler.createRecords() got different http status code = %d, wantCode %d", rec.Code, tt.wantCode)
			}

		})
	}
}
