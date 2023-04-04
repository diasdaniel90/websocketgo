package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	mysqlConfig, err := Envs()
	if err != nil {
		log.Fatal(err)
	}
	db, err := sql.Open("mysql", mysqlConfig.MysqlString())
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	conn, err := connect()
	if err != nil {
		log.Printf("error connecting to websocket: %v", err)
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
				log.Printf("Erro ao decodificar mensagem: %s", err)
				return
			}

			if err := filterMessage(db, payload); err != nil {
				log.Printf("Erro ao filtrar mensagem: %s", err.Error())
				return
			}

			if payload.Status == "waiting" {
				err := saveToDatabaseUsers(db, payload)
				if err != nil {
					log.Printf("error no banco: %s", err)
					return
				}
			}

		case err := <-errChan:
			log.Println(err)
			reconnect(conn, msgChan, errChan)
		}
	}
}
