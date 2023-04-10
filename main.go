package main

import (
	"database/sql"
	"io"
	"log"
	"runtime"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	msgSignalChan := make(chan MsgSignal)

	go listenUDP(msgSignalChan)

	dbConexao, err := sql.Open("mysql", EnvsDatabase())
	if err != nil {
		log.Fatal(err)
	}

	defer dbConexao.Close()

	conn, err := connect()
	if err != nil {
		log.Printf("error connecting to websocket: %v", err)
	}
	defer conn.Close()

	msgChan := make(chan []byte)
	errChan := make(chan error)
	msgStatusChan := make(chan MsgStatus)

	go testeStatus(msgStatusChan, msgSignalChan)

	go readMessages(conn, msgChan, errChan)
	go writePing(conn)
	log.Println("main", runtime.NumGoroutine())

	var wg sync.WaitGroup

	wg.Add(1)

	go teste(&wg, conn, dbConexao, msgChan, errChan, msgStatusChan)

	wg.Wait()
}

type MsgStatus struct {
	IDBet     string `json:"idBet"`
	Timestamp int64  `json:"timestamp"`
	BetStatus string `json:"betStatus"`
	BetColor  int    `json:"betColor"`
	BetRoll   int    `json:"betRoll"`
}

func testeStatus(msgStatusChan <-chan MsgStatus, msgSignalChan <-chan MsgSignal) {
	log.Println("###########11##################")

	mensagens := []MsgSignal{}

	for {
		select {
		case msg, ok := <-msgStatusChan:
			if !ok {
				log.Println("Canal msgStatusChan fechado")

				return
			}

			log.Println("chegou na go func de aposta", msg)
			log.Println("Recebeu sinal msgStatusChan ", mensagens)
		case signalMsg, ok := <-msgSignalChan:
			if !ok {
				log.Println("Canal msgSignalChan fechado")

				return
			}

			mensagens = append(mensagens, signalMsg)
			log.Println("Recebeu sinal msgSignalChan", mensagens)
			log.Println("Recebeu sinal ", signalMsg)
		}
	}
}

func teste(wg *sync.WaitGroup, conn io.Closer, dbConexao *sql.DB, msgChan chan []byte,
	errChan chan error, msgStatusChan chan MsgStatus,
) {
	defer wg.Done()

	for {
		select {
		case msg := <-msgChan:
			payload, err := decodePayload(msg[2:])
			if err != nil {
				log.Printf("Erro ao decodificar mensagem: %s", err)

				return
			}

			Status, err := filterMessage(dbConexao, payload)
			if err != nil {
				log.Printf("Erro ao filtrar mensagem: %s", err.Error())

				return
			}

			if Status != nil {
				msgStatusChan <- *Status
			}

			if payload.Status == "waiting" {
				err := saveToDatabaseUsers(dbConexao, payload)
				if err != nil {
					log.Printf("error no banco: %s", err)

					return
				}
			}

		case err := <-errChan:
			log.Println(err)
			reconnect(conn, msgChan, errChan)
		}
	}
}
