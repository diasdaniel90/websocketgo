package main

import (
	"database/sql"
	"fmt"
	"log"
)

func main() {

	db, err := sql.Open("mysql", "usuario:senha@/nome_do_banco_de_dados")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	conn, err := connect()
	if err != nil {
		log.Fatalf("error connecting to websocket: %v", err)
	}
	defer conn.Close()

	msgChan := make(chan []byte)
	errChan := make(chan error)

	go readMessages(conn, msgChan, errChan)
	go writePing(conn)

	for {
		select {
		case msg := <-msgChan:
			payload, err := decodePayload(msg[2:])
			if err != nil {
				log.Fatalf("Erro ao decodificar mensagem: %s", err)
			} else {
				if err := filterMessage(payload); err != nil {
					log.Fatalf("Erro ao filtrar mensagem: %s", err)
				}
				log.Print(payload)
				if payload.Status == "waiting" {
					err := saveToDatabase(payload)
					if err != nil {
						log.Fatalf("Erro ao inserir mensagem: %s", err)
					}
				}
			}
		case err := <-errChan:
			fmt.Println(err)
			reconnect(conn, msgChan, errChan)
		}
	}
}
