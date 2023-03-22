package main

import (
	"encoding/json"
)

var last_updated_at string = "0"
var last_id string = "0"
var last_id_waiting string = "0"

type Payload struct {
	ID        string `json:"id"`
	Color     int    `json:"color"`
	Roll      int    `json:"roll"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Status    string `json:"status"`
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
