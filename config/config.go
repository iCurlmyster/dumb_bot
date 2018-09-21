package config

import (
	"encoding/json"
	"os"
)

// Config holds values necessary to connect to twitch app
type Config struct {
	OAuth    string `json:"oauth"`
	UserName string `json:"username"`
	Channel  string `json:"channel"`
}

// Configuration opens and pulls out the values from a config.json file
func Configuration(fileName string) (*Config, error) {
	fp, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	var c *Config
	decoder := json.NewDecoder(fp)
	if err := decoder.Decode(&c); err != nil {
		return nil, err
	}
	return c, nil
}
