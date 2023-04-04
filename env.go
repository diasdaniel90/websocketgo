package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	EnvRef        string `json:"envRef"`
	MySQLDatabase string `json:"mySqlDatabase"`
	MySQLUser     string `json:"mySqlUser"`
	MySQLPassword string `json:"mySqlPassword"`
	MySQLHost     string `json:"mySqlHost"`
	MySQLPort     string `json:"mySqlPort"`
}

func Envs() (*Config, error) {
	file, err := os.ReadFile("config.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var cfg Config
	err = json.Unmarshal(file, &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return &cfg, nil
}

func (p Config) MysqlString() string {
	return p.MySQLUser + ":" + p.MySQLPassword + "@tcp(" + p.MySQLHost + ":" + p.MySQLPort + ")/" + p.MySQLDatabase
}
