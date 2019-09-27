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
		userID  string
		secret  string
		fields  fields
		want    []string
		wantErr bool
	}{
		{
			name:   "ok",
			userID: "user",
			secret: "secret",
			want:   []string{"test.com"},
		},
		{
			name:    "authenticate failed",
			userID:  "unsuccess",
			secret:  "unsuccess",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				bufbody := new(bytes.Buffer)
				bufbody.ReadFrom(r.Body)
				if "{\"user_id\":\"user\",\"secret\":\"secret\"}" == bufbody.String() {
					fmt.Fprintln(w, `{"domains": ["test.com"]}`)
				}
			})

			l, _ := net.Listen("tcp", "127.0.0.1:0")
			ts := httptest.Server{
				Listener: l,
				Config:   &http.Server{Handler: handler},
			}
			ts.Start()
			defer ts.Close()

			h := httpAuth{
				Endpoint:      ts.URL,
				RequestHeader: tt.fields.requestHeader,
			}
			got, err := h.Authenticate(tt.userID, tt.secret)
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
