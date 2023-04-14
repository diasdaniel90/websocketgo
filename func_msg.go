package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
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
		// msgStatusChan <- Status

		log.Println("filterMessage Apostas fechadas e resultado", Status)

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

		return &Status, nil
	}

	return nil, err
}
