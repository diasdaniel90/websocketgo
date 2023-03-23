package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

const (
	url          = "wss://api-v2.blaze.com/replication/?EIO=3&transport=websocket"
	pingInterval = 23 * time.Second
	pingMessage  = "2"
	writeWait    = 2 * time.Second
)

func writePing(conn *websocket.Conn) {
	ticker := time.NewTicker(pingInterval)
	defer ticker.Stop()
	for range ticker.C {
		err := conn.WriteMessage(websocket.BinaryMessage, []byte(pingMessage))
		if err != nil {
			log.Printf("error writing ping message: %v", err)
			return
		}
		log.Println("Ping enviado com sucesso")
	}
}

func connect() (*websocket.Conn, error) {
	fmt.Print("connect")
	dialer := websocket.Dialer{
		HandshakeTimeout: 2 * time.Second,
	}

	conn, _, err := dialer.Dial(url, nil)
	if err != nil {
		return nil, fmt.Errorf("error connecting to websocket: %v", err)
	}

	// Envie uma mensagem de assinatura para o servidor
	message := []byte(`420["cmd",{"id":"subscribe","payload":{"room":"double_v2"}}]`)
	err = conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		log.Fatalf("Erro ao enviar mensagem: %s", err)
	}
	log.Println("Assinatura enviada com sucesso")
	return conn, nil
}

func readMessages(conn *websocket.Conn, msgChan chan []byte, errChan chan error) {
	for {
		_, payload, err := conn.ReadMessage()
		if err != nil {
			errChan <- fmt.Errorf("error reading message: %v", err)
			return
		} else {
			if strings.Contains(string(payload), "double.tick") {
				msgChan <- payload
			}

		}

	}
}

func reconnect(conn *websocket.Conn, msgChan chan []byte, errChan chan error) {
	log.Println("Conexão fechada pelo servidor, reconectando...")
	err := conn.Close()
	if err != nil {
		log.Printf("error closing connection: %v", err)
	}
	time.Sleep(2 * time.Second)
	newConn, err := connect()
	if err != nil {
		errChan <- fmt.Errorf("error reconnecting to websocket: %v", err)
		return
	}
	log.Println("Conectado novamente!")
	go readMessages(newConn, msgChan, errChan)
	go writePing(newConn)
}
