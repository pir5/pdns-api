package pdns_api

import (
	"bytes"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func Test_httpAuth_Authenticate(t *testing.T) {
	type fields struct {
		requestHeader map[string]string
	}
	tests := []struct {
		name    string
		token   string
		fields  fields
		want    []string
		wantErr bool
	}{
		{
			name:  "ok",
			token: "ok",
			want:  []string{"test.com"},
		},
		{
			name:    "authenticate failed",
			token:   "unsuccess",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				bufbody := new(bytes.Buffer)
				bufbody.ReadFrom(r.Body)
				if "token=ok" == bufbody.String() {
					fmt.Fprintln(w, `{"domains": ["test.com"]}`)
				}
			})

			l, _ := net.Listen("tcp", "127.0.0.1:8080")
			ts := httptest.Server{
				Listener: l,
				Config:   &http.Server{Handler: handler},
			}
			ts.Start()
			defer ts.Close()

			h := httpAuth{
				endpoint:      "http://localhost:8080",
				requestHeader: tt.fields.requestHeader,
			}
			got, err := h.Authenticate(tt.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("httpAuth.Authenticate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("httpAuth.Authenticate() = %v, want %v", got, tt.want)
			}
		})
	}
}
