package pdns_api

import (
	"github.com/BurntSushi/toml"
)

func NewConfig(confPath string) (Config, error) {
	var conf Config
	defaultConfig(&conf)

	if _, err := toml.DecodeFile(confPath, &conf); err != nil {
		return conf, err
	}

	return conf, nil
}

type Config struct {
	Listen string `toml:"listen"`
}

func defaultConfig(c *Config) {
	c.Listen = "0.0.0.0:8080"
}
