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
				//Verifica se a mensagem é duplicada com base no campo updated_at
				if payload.Status != "waiting" && last_updated_at != payload.UpdatedAt && last_id != payload.ID {
					log.Printf("Mensagem recebida, Enviar msg de socket: %+v", payload)
					last_updated_at = payload.UpdatedAt
					last_id = payload.ID
					//aqui precisa ser enviada uma msg UDP para o servidor
					err := sendUDPMessage("Hello, world!")
					if err != nil {
						// tratar erro de envio
						log.Printf("error sending: %v", err)
					}
				} else if payload.Status == "waiting" && last_id_waiting != payload.ID {
					log.Printf("Mensagem waiting, Enviar msg de socket: %+v", payload)
					last_id_waiting = payload.ID
					//aqui precisa ser enviada uma msg UDP para o servidor
				}

				if payload.Status == "waiting" {
					log.Print("waiting")
				}
			}
		case err := <-errChan:
			fmt.Println(err)
			reconnect(conn, msgChan, errChan)
		}
	}
}
