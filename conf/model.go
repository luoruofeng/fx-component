package conf

import (
	_ "embed"
	"encoding/json"
)

//go:embed config.json
var configStr []byte

func GetConfig() *Config {
	var c *Config
	json.Unmarshal(configStr, &c)
	return c
}

type Config struct {
	Addr         string `json:"addr"`
	Password     string `json:"password"`
	DbNumber     int    `json:"db_number"`
	MaxRetries   int    `json:"max_retries"`
	DialTimeout  int    `json:"dial_timeout"`
	ReadTimeout  int    `json:"read_timeout"`
	WriteTimeout int    `json:"write_timeout"`
}
