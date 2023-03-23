package main

import "log"

func saveToDatabase(payload *Payload) error {
	// Lógica para salvar o payload no banco de dados aqui
	log.Println("se vira para salvar essa parada no banco de dados", payload)
	return nil
}
