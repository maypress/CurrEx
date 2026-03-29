package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Favorites []string `json:"favorites"`
}

func LoadConfig() *Config {
	data, err := os.ReadFile("config.json")
	if err != nil {
		return &Config{}
	}
	
	var cfg Config
	json.Unmarshal(data, &cfg)
	return &cfg
}