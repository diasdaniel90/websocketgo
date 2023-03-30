package main

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	EnvRef        string `json:"EnvRef"`
	MySQLDatabase string `json:"MySQLDatabase"`
	MySQLUser     string `json:"MySQLUser"`
	MySQLPassword string `json:"MySQLPassword"`
	MySQLHost     string `json:"MySQLHost"`
	MySQLPort     string `json:"MySQLPort"`
}

func Envs() (*Config, error) {
	file, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err)
	}

	var cfg Config
	err = json.Unmarshal(file, &cfg)
	if err != nil {
		panic(err)
	}
	return &cfg, nil
}

func (p Config) MysqlString() string {
	return p.MySQLUser + ":" + p.MySQLPassword + "@tcp(" + p.MySQLHost + ":" + p.MySQLPort + ")/" + p.MySQLDatabase
}
