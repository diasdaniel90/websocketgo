package main

import (
	"fmt"
	"log"
)

func main() {
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

				// if payload.Status == "waiting" {
				// 	log.Print("waiting")
				// }
			}
		case err := <-errChan:
			fmt.Println(err)
			reconnect(conn, msgChan, errChan)
		}
	}
}
