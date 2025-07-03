package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	GatewayURL               string `json:"gateway_url"`
	Secret                   string `json:"secret"`
	Zone                     string `json:"zone"`
	HeartbeatIntervalSeconds int    `json:"heartbeat_interval_seconds"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
