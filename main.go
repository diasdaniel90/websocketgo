package main

import (
	"database/sql"
	"io"
	"log"
	"runtime"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const tempoEspera = 4

func main() {
	msgSignalChan := make(chan msgSignalStruct)

	go listenUDP(msgSignalChan)

	dbConexao, err := sql.Open("mysql", envsDatabase())
	if err != nil {
		log.Fatal(err)
	}

	defer dbConexao.Close()

	conn, err := connect()
	if err != nil {
		log.Printf("error connect() to websocket: %v", err)
	}
	defer conn.Close()

	msgChan := make(chan []byte)
	errChan := make(chan error)
	msgStatusChan := make(chan msgStatusStruct)

	go controlBet(msgStatusChan, msgSignalChan)

	go readMessages(conn, msgChan, errChan)
	// go writePing(conn)
	log.Println("main", runtime.NumGoroutine())

	var wg sync.WaitGroup

	wg.Add(1)

	go controlMsg(&wg, conn, dbConexao, msgChan, errChan, msgStatusChan)

	wg.Wait()
}

func validateBet(msgStatusRec msgStatusStruct, sliceBets *[]betBotStruct) {
	log.Print("valida win")

	if msgStatusRec.betStatus != waiting {
		for _, value := range *sliceBets {
			if value.idBet == msgStatusRec.idBet && value.color == msgStatusRec.color {
				value.win = true
				log.Println("vaivai", value.win)
			}
		}

		// log.Println("Color:", msgStatusRec.color)
		// log.Println("resultado", bets)

		*sliceBets = []betBotStruct{}
	}
}

func controlBet(msgStatusChan <-chan msgStatusStruct, msgSignalChan <-chan msgSignalStruct) {
	log.Println("###########11##################")

	sliceSignals := []msgSignalStruct{}
	sliceBets := []betBotStruct{}
	// var valido string
	for {
		select {
		case msgStatusRec, ok := <-msgStatusChan:
			if !ok {
				log.Println("Canal msgStatusChan fechado")

				continue
			}

			if msgStatusRec.betStatus == waiting {
				time.AfterFunc(tempoEspera*time.Second, func() {
					go sinal2Playbet(&sliceSignals, msgStatusRec, &sliceBets)
				})
			}

			validateBet(msgStatusRec, &sliceBets)

		case signalMsg, ok := <-msgSignalChan:
			if !ok {
				log.Println("Canal msgSignalChan fechado")

				continue
			}
			// log.Println("Recebeu sinal msgSignalChan", mensagens)
			// log.Println("Recebeu sinal ", signalMsg)
			sliceSignals = append(sliceSignals, signalMsg)

		default:
			// Faça algo aqui se ambos os canais estiverem vazios.
			// Por exemplo, tente novamente mais tarde.
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func controlMsg(wg *sync.WaitGroup, conn io.Closer, dbConexao *sql.DB, msgChan chan []byte,
	errChan chan error, msgStatusChan chan msgStatusStruct,
) {
	defer wg.Done()

	var lastMsg lastMsgStruct

	for {
		select {
		case msg, ok := <-msgChan:
			if !ok {
				log.Println("Canal msgStatusChan fechado")

				continue
			}

			payload, err := decodePayload(msg[2:])
			if err == nil {
				log.Printf("Erro ao decodificar mensagem: %s", err)

				continue
			}

			Status, err := filterMessage(dbConexao, payload, &lastMsg)
			if err != nil {
				log.Printf("Erro ao filtrar mensagem: %s", err.Error())

				continue
			}

			if payload.Status == "waiting" {
				err := saveToDatabaseUsers(dbConexao, payload)
				if err != nil {
					log.Printf("error no banco: %s", err)

					continue
				}
			}

			if Status != nil {
				msgStatusChan <- *Status
			}

		case err := <-errChan:
			log.Println(err)
			reconnect(conn, msgChan, errChan)
		}
	}
}
