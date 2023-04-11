package main

import (
	"encoding/json"
	"log"
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

func EnvsDatabase() string {
	file, err := os.ReadFile("config.json")
	if err != nil {
		log.Printf("Erro ao ler arquivo: %s", err.Error())
		// return fmt.Errorf("Erro ao ler arquivo: %w", err)

		panic(err.Error())
	}

	var cfg Config

	err = json.Unmarshal(file, &cfg)
	if err != nil {
		log.Printf("Erro ao ler arquivo: %s", err.Error())
		panic(err.Error())
	}

	return cfg.MySQLUser + ":" + cfg.MySQLPassword +
		"@tcp(" + cfg.MySQLHost + ":" + cfg.MySQLPort + ")/" + cfg.MySQLDatabase
}
