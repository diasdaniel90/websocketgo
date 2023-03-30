package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"time"
)

var lastUpdatedAt string = "0"
var lastId string = "0"
var lastIdWaiting string = "0"

const layout = "2006-01-02T15:04:05.000Z"

type Payload struct {
	IdBet                string `json:"id"`
	Color                int    `json:"color"`
	Roll                 int    `json:"roll"`
	CreatedAt            string `json:"created_at"`
	Timestamp            int64
	UpdatedAt            string  `json:"updated_at"`
	Status               string  `json:"status"`
	TotalRedEurBet       float64 `json:"total_red_eur_bet"`
	TotalRedBetsPlaced   int     `json:"total_red_bets_placed"`
	TotalWhiteEurBet     float64 `json:"total_white_eur_bet"`
	TotalWhiteBetsPlaced int     `json:"total_white_bets_placed"`
	TotalBlackEurBet     float64 `json:"total_black_eur_bet"`
	TotalBlackBetsPlaced int     `json:"total_black_bets_placed"`
	TotalBetsPlaced      int
	TotalEurBet          float64
	TotalRetentionEur    float64
	Bets                 []Bet `json:"bets"`
}

func (p *Payload) calculateTotalBetsPlaced() {
	p.TotalBetsPlaced = int(p.TotalRedBetsPlaced + p.TotalWhiteBetsPlaced + p.TotalBlackBetsPlaced)
}

func (p *Payload) calculateTotalBetsEur() {
	p.TotalEurBet = float64(p.TotalRedEurBet + p.TotalWhiteEurBet + p.TotalBlackEurBet)
}

func (p *Payload) calculateTotalRetentionEur() {
	switch p.Color {
	case 1:
		p.TotalRetentionEur = float64(p.TotalEurBet - p.TotalRedEurBet*2)
	case 2:
		p.TotalRetentionEur = float64(p.TotalEurBet - p.TotalBlackEurBet*2)
	case 0:
		p.TotalRetentionEur = float64(p.TotalEurBet - p.TotalWhiteEurBet*14)
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
	//payload := data[1].(map[string]interface{})["payload"]
	var data []json.RawMessage
	if err := json.Unmarshal([]byte(message), &data); err != nil {
		panic(err)
	}

	var payload Payload
	if err := json.Unmarshal(data[1], &struct{ Payload *Payload }{&payload}); err != nil {
		panic(err)
	}
	// Retorna a mensagem decodificada
	return &payload, nil
}

func filterMessage(db *sql.DB, payload *Payload) error {
	//Verifica se a mensagem é duplicada com base no campo updated_at
	if payload.Status != "waiting" && lastUpdatedAt != payload.UpdatedAt && lastId != payload.IdBet {
		lastUpdatedAt = payload.UpdatedAt
		lastId = payload.IdBet
		tComplete, _ := time.Parse(layout, payload.CreatedAt)
		payload.Timestamp = tComplete.Unix()

		payload.calculateTotalBetsPlaced()
		payload.calculateTotalBetsEur()
		payload.calculateTotalRetentionEur()
		err := saveToDatabase(db, payload)
		if err != nil {
			log.Printf("error no banco: %v", err)
			return err
		}
		err = sendUDPMessage(payload)
		if err != nil {
			log.Printf("error sending: %v", err)
			return err
		}
		log.Println("Apostas fechadas e resultado")
	} else if payload.Status == "waiting" && lastIdWaiting != payload.IdBet {
		lastIdWaiting = payload.IdBet
		tWaiting, _ := time.Parse(layout, payload.CreatedAt)
		payload.Timestamp = tWaiting.Unix()
		err := sendUDPMessage(payload)
		if err != nil {
			log.Printf("error sending: %v", err)
			return err
		}
		log.Println("Pronto para apostar")

	}
	return nil
}
