package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"runtime"
	"time"
)

var (
	lastUpdatedAt = "revive"
	lastID        = "revive"
	lastIDWaiting = "revive"
)

const (
	layout     = "2006-01-02T15:04:05.000Z"
	waiting    = "waiting"
	white      = 0
	red        = 1
	black      = 2
	FatorWhite = 14
	FatorColor = 2
)

type Payload struct {
	IDBet                string  `json:"id"`
	Color                int     `json:"color"`
	Roll                 int     `json:"roll"`
	CreatedAt            string  `json:"created_at"`
	Timestamp            int64   `json:"timestamp"`
	UpdatedAt            string  `json:"updated_at"`
	Status               string  `json:"status"`
	TotalRedEurBet       float64 `json:"total_red_eur_bet"`
	TotalRedBetsPlaced   int     `json:"total_red_bets_placed"`
	TotalWhiteEurBet     float64 `json:"total_white_eur_bet"`
	TotalWhiteBetsPlaced int     `json:"total_white_bets_placed"`
	TotalBlackEurBet     float64 `json:"total_black_eur_bet"`
	TotalBlackBetsPlaced int     `json:"total_black_bets_placed"`
	TotalBetsPlaced      int     `json:"totalBetsPlaced"`
	TotalEurBet          float64 `json:"totalEurBet"`
	TotalRetentionEur    float64 `json:"totalRetentionEur"`
	Bets                 []Bet   `json:"bets"`
}

func (p *Payload) calculateTotalBetsPlaced() {
	p.TotalBetsPlaced = p.TotalRedBetsPlaced + p.TotalWhiteBetsPlaced + p.TotalBlackBetsPlaced
}

func (p *Payload) calculateTotalBetsEur() {
	p.TotalEurBet = p.TotalRedEurBet + p.TotalWhiteEurBet + p.TotalBlackEurBet
}

func (p *Payload) calculateTotalRetentionEur() {
	switch p.Color {
	case red:
		p.TotalRetentionEur = p.TotalEurBet - p.TotalRedEurBet*FatorColor
	case black:
		p.TotalRetentionEur = p.TotalEurBet - p.TotalBlackEurBet*FatorColor
	case white:
		p.TotalRetentionEur = p.TotalEurBet - p.TotalWhiteEurBet*FatorWhite
	}
}

type Bet struct {
	IDBetUser    string  `json:"id"`
	Color        int     `json:"color"`
	Amount       float32 `json:"amount"`
	CurrencyType string  `json:"currency_type"`
	Status       string  `json:"status"`
	User         struct {
		IDStr string `json:"id_str"`
	} `json:"user"`
}

func decodePayload(message []byte) (*Payload, error) {
	log.Println("Gotoutine", runtime.NumGoroutine())

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

func filterMessage(dbConexao *sql.DB, payload *Payload) (*MsgStatus, error) {
	// Verifica se a mensagem é duplicada com base no campo updated_at
	var err error

	if payload.Status != waiting && lastUpdatedAt != payload.UpdatedAt && lastID != payload.IDBet {
		lastUpdatedAt = payload.UpdatedAt
		lastID = payload.IDBet
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
			BetColor:  payload.Color,
			BetRoll:   payload.Roll,
		}
		// msgStatusChan <- Status
		log.Println("#############################")

		// if err := sendUDPMessage(payload); err != nil {
		// 	fmt.Errorf("error saveToDatabase: %w", err)
		// }

		log.Println("filterMessage Apostas fechadas e resultado")

		return &Status, nil
	} else if payload.Status == waiting && lastIDWaiting != payload.IDBet {
		lastIDWaiting = payload.IDBet
		tWaiting, _ := time.Parse(layout, payload.CreatedAt)
		payload.Timestamp = tWaiting.Unix()
		// err := sendUDPMessage(payload)
		// if err != nil {
		// 	return fmt.Errorf("error connecting to websocket: %w", err)
		// }
		log.Println("filterMessage Pronto para apostar")

		Status := MsgStatus{
			IDBet:     payload.IDBet,
			Timestamp: payload.Timestamp,
			BetStatus: payload.Status,
			BetColor:  payload.Color,
			BetRoll:   payload.Roll,
		}
		// msgStatusChan <- Status
		log.Println("#############################")

		return &Status, nil
	}

	return nil, err
}
