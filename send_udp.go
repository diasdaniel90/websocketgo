package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
)

const address = "127.0.0.1:20001"

//({'ID_bet': "B", 'timestamp': timestamp_status, 'bet_status': 'waiting'})
//({'ID_bet': "B", 'timestamp': timestamp_status, 'bet_status': 'rolling', 'bet_color': 2, 'bet_roll': 13})

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
		return err
	}
	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	//message := Mensagem{payload.IdBet, payload.Timestamp, payload.Status, payload.Color, payload.Roll}

	message := Mensagem{
		ID_bet:     payload.IdBet,
		Timestamp:  payload.Timestamp,
		Bet_status: payload.Status,
		Bet_color:  payload.Color,
		Bet_roll:   payload.Roll,
	}

	messageJson, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Erro ao converter para JSON:", err)
		return err
	}

	_, err = conn.Write([]byte(messageJson))
	if err != nil {
		return err
	}
	//fmt.Print(n)
	//fmt.Printf("message enviada : %+v\n", messageJson)
	log.Print("message enviada : ", messageJson)
	return nil
}
