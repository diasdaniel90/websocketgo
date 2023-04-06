package main

import (
	"fmt"
	"io"
	"log"
	"runtime"
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

func connect() (*websocket.Conn, error) {
	log.Print("connect")

	dialer := websocket.Dialer{
		HandshakeTimeout: writeWait,
	}

	conn, _, err := dialer.Dial(url, nil)
	if err != nil {
		return nil, fmt.Errorf("error connecting to websocket: %w", err)
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

func reconnect(conn io.Closer, msgChan chan []byte, errChan chan error) {
	log.Println("Conexão fechada pelo servidor, reconectando...")

	if err := conn.Close(); err != nil {
		log.Printf("error closing connection: %v", err)
	}

	time.Sleep(writeWait)

	newConn, err := connect()
	if err != nil {
		errChan <- fmt.Errorf("error reconnecting to websocket: %w", err)

		return
	}

	log.Println("Conectado novamente!")
	log.Println("reconect", runtime.NumGoroutine())

	go readMessages(newConn, msgChan, errChan)
	go writePing(newConn)
	log.Println("reconect2", runtime.NumGoroutine())
}

func readMessages(conn *websocket.Conn, msgChan chan []byte, errChan chan error) {
	for {
		_, payload, err := conn.ReadMessage()
		if err != nil {
			errChan <- fmt.Errorf("error reading message: %w", err)

			return
		} else if strings.Contains(string(payload), "double.tick") {
			msgChan <- payload
		}
	}
}

func writePing(conn *websocket.Conn) {
	log.Println("writePing")

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
