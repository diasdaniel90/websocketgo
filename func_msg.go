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

func (p *Payload) calculateTotalBetsPlaced() {
	p.TotalBetsPlaced = p.TotalRedBetsPlaced + p.TotalWhiteBetsPlaced + p.TotalBlackBetsPlaced
}

func (p *Payload) calculateTotalBetsEur() {
	p.TotalEurBet = p.TotalRedEurBet + p.TotalWhiteEurBet + p.TotalBlackEurBet
}

func (p *Payload) calculateTotalRetentionEur() {
	switch p.Color {
	case red:
		p.TotalRetentionEur = p.TotalEurBet - p.TotalRedEurBet*fatorColor
	case black:
		p.TotalRetentionEur = p.TotalEurBet - p.TotalBlackEurBet*fatorColor
	case white:
		p.TotalRetentionEur = p.TotalEurBet - p.TotalWhiteEurBet*fatorWhite
	}
}

func decodePayload(message []byte) (*Payload, error) {
	// log.Println("Gotoutine", runtime.NumGoroutine())
	var data []json.RawMessage
	if err := json.Unmarshal(message, &data); err != nil {
		return nil, fmt.Errorf("error unmarshaling payload:: %w", err)
	}

	var payload Payload
	if err := json.Unmarshal(data[1], &struct {
		Payload *Payload `json:"payload"`
	}{&payload}); err != nil {
		return nil, fmt.Errorf("error unmarshaling payload:: %w", err)
	}
	// Retorna a mensagem decodificada
	return &payload, nil
}

func filterMessage(dbConexao *sql.DB, payload *Payload, lastMsg *LastMsg) (*MsgStatus, error) {
	// Verifica se a mensagem é duplicada com base no campo updated_at
	var err error

	if payload.Status != waiting && lastMsg.LastUpdatedAt != payload.UpdatedAt &&
		lastMsg.LastID != payload.IDBet && lastMsg.LastIDWaiting == payload.IDBet {
		lastMsg.LastUpdatedAt = payload.UpdatedAt
		lastMsg.LastID = payload.IDBet
		tComplete, _ := time.Parse(layout, payload.CreatedAt)
		payload.Timestamp = tComplete.Unix()

		payload.calculateTotalBetsPlaced()
		payload.calculateTotalBetsEur()
		payload.calculateTotalRetentionEur()

		if err = saveToDatabase(dbConexao, payload); err != nil {
			return nil, fmt.Errorf("error saveToDatabase: %w", err)
		}

		Status := MsgStatus{
			IDBet:     payload.IDBet,
			Timestamp: payload.Timestamp,
			BetStatus: payload.Status,
			Color:     payload.Color,
			BetRoll:   payload.Roll,
		}
		// msgStatusChan <- Status

		log.Println("filterMessage Apostas fechadas e resultado")

		return &Status, nil
	} else if payload.Status == waiting && lastMsg.LastIDWaiting != payload.IDBet {
		lastMsg.LastIDWaiting = payload.IDBet
		tWaiting, _ := time.Parse(layout, payload.CreatedAt)
		payload.Timestamp = tWaiting.Unix()

		Status := MsgStatus{
			IDBet:     payload.IDBet,
			Timestamp: payload.Timestamp,
			BetStatus: payload.Status,
			Color:     payload.Color,
			BetRoll:   payload.Roll,
		}
		log.Println("filterMessage pronto para apostar")

		return &Status, nil
	}

	return nil, err
}
