package controller

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
		recordModel model.RecordModeler
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
			wantCode: http.StatusOK,
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
		recordModel model.RecordModeler
		domainModel model.DomainModeler
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

			if err := h.updateRecord(ctx); (err != nil) != tt.wantErr {
				t.Errorf("recordHandler.updateRecord() error = %v, wantErr %v", err, tt.wantErr)
			}

			if rec.Code != tt.wantCode {
				t.Errorf("recordHandler.updateRecords() got different http status code = %d, wantCode %d", rec.Code, tt.wantCode)
			}

		})
	}
}

func Test_recordHandler_createRecord(t *testing.T) {
	type fields struct {
		recordModel model.RecordModeler
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
				recordModel: &recordModelStub{},
			},
			wantErr:  false,
			wantCode: http.StatusCreated,
			queryID:  "1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &recordHandler{
				recordModel: tt.fields.recordModel,
				domainModel: tt.fields.domainModel,
			}
			ctx, rec := dummyContext(t, "POST", "/records", nil)

			if err := h.createRecord(ctx); (err != nil) != tt.wantErr {
				t.Errorf("recordHandler.createRecord() error = %v, wantErr %v", err, tt.wantErr)
			}

			if rec.Code != tt.wantCode {
				t.Errorf("recordHandler.createRecords() got different http status code = %d, wantCode %d", rec.Code, tt.wantCode)
			}

		})
	}
}

func Test_recordHandler_enableRecord(t *testing.T) {
	type fields struct {
		recordModel model.RecordModeler
		domainModel model.DomainModeler
	}
	tests := []struct {
		name     string
		fields   fields
		wantCode int
		wantErr  bool
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
			name: "not found",
			fields: fields{
				domainModel: &domainModelStub{},
				recordModel: &recordModelStub{},
			},
			wantErr:  false,
			wantCode: http.StatusNotFound,
			queryID:  "0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &recordHandler{
				recordModel: tt.fields.recordModel,
				domainModel: tt.fields.domainModel,
			}
			ctx, rec := dummyContext(t, "PUT", "/records/enable/:", nil)
			ctx.SetParamNames("id")
			ctx.SetParamValues(tt.queryID)
			if err := h.enableRecord(ctx); (err != nil) != tt.wantErr {
				t.Errorf("recordHandler.enableRecord() error = %v, wantErr %v", err, tt.wantErr)
			}
			if rec.Code != tt.wantCode {
				t.Errorf("recordHandler.enableRecords() got different http status code = %d, wantCode %d", rec.Code, tt.wantCode)
			}
		})
	}
}

func Test_recordHandler_disableRecord(t *testing.T) {
	type fields struct {
		recordModel model.RecordModeler
		domainModel model.DomainModeler
	}
	tests := []struct {
		name     string
		fields   fields
		wantCode int
		wantErr  bool
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
			name: "not found",
			fields: fields{
				domainModel: &domainModelStub{},
				recordModel: &recordModelStub{},
			},
			wantErr:  false,
			wantCode: http.StatusNotFound,
			queryID:  "0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &recordHandler{
				recordModel: tt.fields.recordModel,
				domainModel: tt.fields.domainModel,
			}
			ctx, rec := dummyContext(t, "PUT", "/records/disable/:", nil)
			ctx.SetParamNames("id")
			ctx.SetParamValues(tt.queryID)
			if err := h.disableRecord(ctx); (err != nil) != tt.wantErr {
				t.Errorf("recordHandler.disableRecord() error = %v, wantErr %v", err, tt.wantErr)
			}
			if rec.Code != tt.wantCode {
				t.Errorf("recordHandler.disableRecords() got different http status code = %d, wantCode %d", rec.Code, tt.wantCode)
			}
		})
	}
}
