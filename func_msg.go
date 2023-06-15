package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"runtime"
	"sync"
	"time"
)

const (
	layout     = "2006-01-02T15:04:05.000Z"
	waiting    = "waiting"
	white      = 0
	red        = 1
	black      = 2
	fatorWhite = 14
	fatorColor = 2
)

func (p *payloadStruct) calculateTotals() {
	p.TotalBetsPlaced = p.TotalRedBetsPlaced + p.TotalWhiteBetsPlaced + p.TotalBlackBetsPlaced

	p.TotalEurBet = p.TotalRedEurBet + p.TotalWhiteEurBet + p.TotalBlackEurBet

	switch p.Color {
	case red:
		p.TotalRetentionEur = p.TotalEurBet - p.TotalRedEurBet*fatorColor
	case black:
		p.TotalRetentionEur = p.TotalEurBet - p.TotalBlackEurBet*fatorColor
	case white:
		p.TotalRetentionEur = p.TotalEurBet - p.TotalWhiteEurBet*fatorWhite
	}
}

func decodePayload(message []byte) (*payloadStruct, error) {
	// log.Println("Gotoutine", runtime.NumGoroutine())
	var data []json.RawMessage
	if err := json.Unmarshal(message, &data); err != nil {
		return nil, fmt.Errorf("error unmarshaling payload:: %w", err)
	}

	var payload payloadStruct
	if err := json.Unmarshal(data[1], &struct {
		Payload *payloadStruct `json:"payload"`
	}{&payload}); err != nil {
		return nil, fmt.Errorf("error unmarshaling payload:: %w", err)
	}

	// Retorna a mensagem decodificada
	return &payload, nil
}

func filterMessage(dbConexao *sql.DB, payload *payloadStruct, lastMsg *lastMsgStruct) (*msgStatusStruct, error) {
	// Verifica se a mensagem é duplicada com base no campo updated_at
	var err error

	if payload.Status != waiting && lastMsg.lastUpdatedAt != payload.UpdatedAt &&
		lastMsg.lastID != payload.IDBet && lastMsg.lastIDWaiting == payload.IDBet {
		lastMsg.lastUpdatedAt = payload.UpdatedAt
		lastMsg.lastID = payload.IDBet
		tComplete, _ := time.Parse(layout, payload.CreatedAt)
		payload.Timestamp = tComplete.Unix()

		payload.calculateTotals()

		if err = saveToDatabase(dbConexao, payload); err != nil {
			return nil, fmt.Errorf("error saveToDatabase: %w", err)
		}

		Status := msgStatusStruct{
			idBet:     payload.IDBet,
			timestamp: payload.Timestamp,
			betStatus: payload.Status,
			color:     payload.Color,
			betRoll:   payload.Roll,
		}

		log.Println("filterMessage Apostas fechadas e resultado", Status)
		log.Println("runtime.NumGoroutine()", runtime.NumGoroutine())

		return &Status, nil
	}

	if payload.Status == waiting && lastMsg.lastIDWaiting != payload.IDBet {
		lastMsg.lastIDWaiting = payload.IDBet
		tWaiting, _ := time.Parse(layout, payload.CreatedAt)
		payload.Timestamp = tWaiting.Unix()

		Status := msgStatusStruct{
			idBet:     payload.IDBet,
			timestamp: payload.Timestamp,
			betStatus: payload.Status,
			color:     payload.Color,
			betRoll:   payload.Roll,
		}
		log.Println("filterMessage pronto para apostar", Status)
		log.Println("runtime.NumGoroutine()", runtime.NumGoroutine())

		return &Status, nil
	}

	return nil, err
}

func controlMsg(wg *sync.WaitGroup, conn io.Closer, dbConexao *sql.DB, msgChanWebsocket chan []byte,
	errChan chan error, msgStatusChan chan msgStatusStruct, msgSignalChan chan msgSignalStruct,
) {
	defer wg.Done()

	var lastMsg lastMsgStruct

	var rankBet betsUsersStruct

	for {
		select {
		case msg, ok := <-msgChanWebsocket:
			if !ok {
				log.Println("Canal msgChanWebsocket fechado")

				continue
			}

			payload, err := decodePayload(msg[2:])
			if err != nil {
				log.Printf("Erro ao decodificar mensagem: %s", err)

				continue
			}

			Status, err := filterMessage(dbConexao, payload, &lastMsg)
			if err != nil {
				log.Printf("Erro ao filtrar mensagem: %s", err.Error())

				continue
			}

			if payload.Status == "waiting" {
				for _, bet := range payload.Bets {
					if bet.Amount > rankBet.Amount && bet.Color != 0 {
						rankBet = bet
					}
				}

				go saveToDatabaseUsers(dbConexao, *payload)
			}

			if Status != nil {
				msgStatusChan <- *Status
				if Status.betStatus == "waiting" {
					time.AfterFunc(10*time.Second, func() {
						go sinalRank(&rankBet, msgSignalChan, payload.IDBet)
					})
				}

				if Status.betStatus == "complete" || Status.betStatus == "rolling" {
					fmt.Println("#######MELHOR APOSTA#######", rankBet, payload.IDBet)
					rankBet = betsUsersStruct{}
				}

			}

		case err := <-errChan:
			log.Println(err)
			reconnect(conn, msgChanWebsocket, errChan)
		}
	}
}

func sinalRank(rank *betsUsersStruct, msgSignalChan chan msgSignalStruct, idbet string) {
	fmt.Println("#######APOSTA QUE DEU PARA PEGAR#######", rank, idbet)
	msgSignal := msgSignalStruct{
		Type:      "realtime",
		Timestamp: 0.0,
		Color:     rank.Color,
		Source:    1,
	}

	msgSignalChan <- msgSignal
}
