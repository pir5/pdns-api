package pdns_api

import (
	"github.com/BurntSushi/toml"
)

const AuthTypeHTTP = "http"
const AuthTypeToken = "token"

func NewConfig(confPath string) (Config, error) {
	var conf Config
	defaultConfig(&conf)

	if _, err := toml.DecodeFile(confPath, &conf); err != nil {
		return conf, err
	}

	return conf, nil
}

type Config struct {
	Listen    string `toml:"listen"`
	TokenAuth tokenAuth
	DB        database `toml:"database"`
}

type database struct {
	Host     string
	Port     int
	DBName   string `toml:"dbname"`
	UserName string `toml:"username"`
	Password string
}

func defaultConfig(c *Config) {
	c.Listen = "0.0.0.0:8080"
	c.DB.Host = "localhost"
	c.DB.Port = 3306
	c.DB.UserName = "root"
	c.DB.DBName = "pdns"
}

type tokenAuth struct {
	AuthType string // token or http
	Tokens   []string
	HttpAuth httpAuth
}

func (c Config) IsTokenAuth() bool {
	return c.TokenAuth.AuthType == AuthTypeToken
}

func (c Config) IsHTTPAuth() bool {
	return c.TokenAuth.AuthType == AuthTypeHTTP
}
