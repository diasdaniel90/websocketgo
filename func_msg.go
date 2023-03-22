package main

import (
	"encoding/json"
)

var last_updated_at string = "0"
var last_id string = "0"
var last_id_waiting string = "0"

type Payload struct {
	IdBet                string  `json:"id"`
	Color                int     `json:"color"`
	Roll                 int     `json:"roll"`
	CreatedAt            string  `json:"created_at"`
	UpdatedAt            string  `json:"updated_at"`
	Status               string  `json:"status"`
	TotalRedEurBet       float64 `json:"total_red_eur_bet"`
	TotalRedBetsPlaced   int     `json:"total_red_bets_placed"`
	TotalWhiteEurBet     float64 `json:"total_white_eur_bet"`
	TotalWhiteBetsPlaced int     `json:"total_white_bets_placed"`
	TotalBlackEurBet     float64 `json:"total_black_eur_bet"`
	TotalBlackBetsPlaced int     `json:"total_black_bets_placed"`
	Bets                 []Bet   `json:"bets"`
}

type Bet struct {
	IDBets       string  `json:"id"`
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
