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

	var auxEurBet totalEurBetStruct

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
				// log.Printf("waiting red: %f , black:%f\n", payload.TotalRedEurBet, payload.TotalBlackEurBet)
				for _, bet := range payload.Bets {
					if bet.Amount > rankBet.Amount && bet.Color != 0 {
						rankBet = bet
					}
				}
				auxEurBet = totalEurBetStruct{
					TotalRedEurBet:   payload.TotalRedEurBet,
					TotalBlackEurBet: payload.TotalBlackEurBet,
				}
				go saveToDatabaseUsers(dbConexao, *payload)
			}

			if Status != nil {
				msgStatusChan <- *Status
				if Status.betStatus == "waiting" {
					time.AfterFunc(12*time.Second, func() {
						log.Printf("after red: %f , black:%f\n", auxEurBet.TotalRedEurBet, auxEurBet.TotalBlackEurBet)
						sinalLogico(auxEurBet.verifySmallbetEur(), rankBet.Color, msgSignalChan, 3)
					})
				}

				if Status.betStatus == "complete" || Status.betStatus == "rolling" {
					log.Printf("final red: %f , black:%f\n", payload.TotalRedEurBet, payload.TotalBlackEurBet)
					log.Println("final COR COM MENOR VALOR:", payload.verifySmallbet())
				}

			}

		case err := <-errChan:
			log.Println(err)
			reconnect(conn, msgChanWebsocket, errChan)
		}
	}
}

func sinalLogico(colorSmallEur int, rankColor int, msgSignalChan chan msgSignalStruct, source int) {
	// log.Printf("sinalLogico red: %f , black:%f\n", auxEurBet.TotalRedEurBet, auxEurBet.TotalBlackEurBet)
	// log.Println("sinalLogico COR COM MENOR VALOR:", auxEurBet.verifySmallbetEur())
	if colorSmallEur == rankColor {
		msgSignal := msgSignalStruct{
			Type:      "realtime",
			Timestamp: 0.0,
			Color:     colorSmallEur,
			Source:    source,
		}

		msgSignalChan <- msgSignal
	}
}
