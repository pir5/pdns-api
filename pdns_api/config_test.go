package pdns_api

import (
	"os"
	"reflect"
	"testing"
)

func TestNewConfig(t *testing.T) {
	type args struct {
		confPath string
	}
	tests := []struct {
		name    string
		args    args
		want    *Config
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				confPath: "./test.toml",
			},
			want: &Config{
				Listen: "0.0.0.0:1000",
				DB: database{
					Host:     "127.0.0.1",
					Port:     3306,
					DBName:   "testdb",
					UserName: "testuser",
					Password: "testpassword",
				},
			},
		},
	}
	for _, tt := range tests {
		os.Setenv("PIR5_DATABASE_PASSWORD", "testpassword")
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewConfig(tt.args.confPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
