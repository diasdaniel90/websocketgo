package main

import (
	"encoding/json"
	"log"
	"os"
)

func envsDatabase() string {
	file, err := os.ReadFile("config.json")
	if err != nil {
		log.Printf("Erro ao ler arquivo: %s", err.Error())
		// return fmt.Errorf("Erro ao ler arquivo: %w", err)

		panic(err.Error())
	}

	var cfg configStruct

	err = json.Unmarshal(file, &cfg)
	if err != nil {
		log.Printf("Erro ao ler arquivo: %s", err.Error())
		panic(err.Error())
	}

	return cfg.MySQLUser + ":" + cfg.MySQLPassword +
		"@tcp(" + cfg.MySQLHost + ":" + cfg.MySQLPort + ")/" + cfg.MySQLDatabase
}
