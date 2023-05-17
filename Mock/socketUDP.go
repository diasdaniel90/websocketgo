package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
)

const address = "127.0.0.1:1234"

type msgSignal struct {
	Type      string  `json:"type"`
	Timestamp float64 `json:"time"`
	Color     int     `json:"color"`
	Source    int     `json:"source"`
}

func main() {
	message := &msgSignal{
		Type:      "realtime",
		Timestamp: 111,
		Color:     1,
		Source:    111,
	}

	message2 := &msgSignal{
		Type:      "realtime",
		Timestamp: 222,
		Color:     2,
		Source:    1111,
	}

	if err := sendUDPMessage(message); err != nil {
		log.Printf("Erro ao decodificar mensagem: %s", err)
	}

	if err := sendUDPMessage(message2); err != nil {
		log.Printf("Erro ao decodificar mensagem: %s", err)
	}
}

func sendUDPMessage(message *msgSignal) error {
	log.Println(message)

	serverAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return fmt.Errorf("error connecting to udp: %w", err)
	}

	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		return fmt.Errorf("error send udp: %w", err)
	}

	defer conn.Close()

	messageJSON, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("error  json.Marshal(message): %w", err)
	}

	_, err = conn.Write(messageJSON)
	if err != nil {
		return fmt.Errorf("error send udp: %w", err)
	}

	return nil
}
