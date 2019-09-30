package pdns_api

import (
	"reflect"
	"strings"

	"github.com/spf13/viper"
)

const AuthTypeHTTP = "http"
const AuthTypeToken = "token"

func BindEnvs(iface interface{}, parts ...string) {
	ifv := reflect.ValueOf(iface)
	ift := reflect.TypeOf(iface)
	for i := 0; i < ift.NumField(); i++ {
		v := ifv.Field(i)
		t := ift.Field(i)
		tv, ok := t.Tag.Lookup("mapstructure")
		if !ok {
			continue
		}
		switch v.Kind() {
		case reflect.Struct:
			BindEnvs(v.Interface(), append(parts, tv)...)
		default:
			viper.BindEnv(strings.Join(append(parts, tv), "."))
		}
	}
}
func NewConfig(confPath string) (Config, error) {
	var conf Config
	defaultConfig(&conf)
	viper.SetConfigFile(confPath)

	viper.AutomaticEnv()
	viper.SetEnvPrefix("PIR5")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		return conf, err
	}

	BindEnvs(conf)
	if err := viper.Unmarshal(&conf); err != nil {
		return conf, err
	}

	return conf, nil
}

type Config struct {
	Listen   string   `mapstructure:"listen"`
	Endpoint string   `mapstructure:"endpoint"`
	Auth     auth     `mapstructure:"auth"`
	DB       database `mapstructure:"database"`
}

type database struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	DBName   string `mapstructure:"dbname"`
	UserName string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

func defaultConfig(c *Config) {
	c.Listen = "0.0.0.0:8080"
	c.DB.Host = "localhost"
	c.DB.Port = 3306
	c.DB.UserName = "root"
	c.DB.DBName = "pdns"
}

type auth struct {
	AuthType string   `mapstructure:"auth_type"` // token or http
	Tokens   []string `mapstructure:"tokens"`
	HttpAuth httpAuth `mapstructure:"http_auth"`
}

func (c Config) IsTokenAuth() bool {
	return c.Auth.AuthType == AuthTypeToken
}

func (c Config) IsHTTPAuth() bool {
	return c.Auth.AuthType == AuthTypeHTTP
}
