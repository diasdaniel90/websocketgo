package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	mysqlConfig, err := Envs()

	db, err := sql.Open("mysql", mysqlConfig.MysqlString())
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
				if err := filterMessage(db, payload); err != nil {
					log.Fatalf("Erro ao filtrar mensagem: %s", err)
				}
				if payload.Status == "waiting" {
					err := saveToDatabaseUsers(db, payload)
					if err != nil {
						log.Fatalf("error no banco: %s", err)
					}
				}
			}
		case err := <-errChan:
			fmt.Println(err)
			reconnect(conn, msgChan, errChan)
		}
	}
}
