package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
)

const address = "127.0.0.1:20001"

type Mensagem struct {
	ID_bet     string `json:"ID_bet"`
	Timestamp  int64  `json:"timestamp"`
	Bet_status string `json:"bet_status"`
	Bet_color  int    `json:"bet_color"`
	Bet_roll   int    `json:"bet_roll"`
}

func sendUDPMessage(payload *Payload) error {
	serverAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return fmt.Errorf("error connecting to websocket: %w", err)
	}
	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		log.Printf("error writing ping message: %v", err)
		return err
	}
	defer conn.Close()

	message := Mensagem{
		ID_bet:     payload.IdBet,
		Timestamp:  payload.Timestamp,
		Bet_status: payload.Status,
		Bet_color:  payload.Color,
		Bet_roll:   payload.Roll,
	}

	messageJson, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("error connecting to websocket: %w", err)
	}

	_, err = conn.Write([]byte(messageJson))
	if err != nil {
		return fmt.Errorf("error connecting to websocket: %w", err)
	}

	return nil
}
