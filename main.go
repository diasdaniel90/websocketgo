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

func controlBet(msgStatusChan <-chan msgStatusStruct, msgSignalChan <-chan msgSignalStruct) {
	log.Println("###########11##################")

	mensagens := []msgSignalStruct{}
	bets := []betBotStruct{}
	// var valido string
	for {
		select {
		case msgStatusRec, ok := <-msgStatusChan:
			if !ok {
				log.Println("Canal msgStatusChan fechado")

				return
			}

			if msgStatusRec.betStatus == waiting {
				time.AfterFunc(tempoEspera*time.Second, func() {
					go sinal2Playbet(&mensagens, msgStatusRec, &bets)
				})
			}

			if msgStatusRec.betStatus != waiting {
				for i := range bets {
					log.Println("vai", bets[i].idBet, msgStatusRec.idBet, bets[i].color, msgStatusRec.color)

					if bets[i].idBet == msgStatusRec.idBet && bets[i].color == msgStatusRec.color {
						bets[i].win = true
						log.Println("vaivai", bets[i].win)
					}
				}

				// log.Println("Color:", msgStatusRec.color)

				// log.Println("resultado", bets)

				bets = []betBotStruct{}
			}

		case signalMsg, ok := <-msgSignalChan:
			if !ok {
				log.Println("Canal msgSignalChan fechado")

				return
			}
			// log.Println("Recebeu sinal msgSignalChan", mensagens)
			// log.Println("Recebeu sinal ", signalMsg)
			mensagens = append(mensagens, signalMsg)

		default:
			// Faça algo aqui se ambos os canais estiverem vazios.
			// Por exemplo, tente novamente mais tarde.
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func sinal2Playbet(signals *[]msgSignalStruct, msgStatus msgStatusStruct, bets *[]betBotStruct) {
	log.Println("Executando a função após 4 segundos...", signals, msgStatus)

	for _, value := range *signals {
		bet := betBotStruct{
			idBet:     msgStatus.idBet,
			timestamp: msgStatus.timestamp,
			color:     value.Color,
			source:    value.Source,
			win:       false,
			status:    false,
		}
		*bets = append(*bets, bet)

		log.Println("value", value)
	}

	*signals = []msgSignalStruct{}

	log.Println("apostas feitas", bets)
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

				return
			}

			payload, err := decodePayload(msg[2:])
			if err == nil {
				log.Printf("Erro ao decodificar mensagem: %s", err)

				return
			}

			Status, err := filterMessage(dbConexao, payload, &lastMsg)
			if err != nil {
				log.Printf("Erro ao filtrar mensagem: %s", err.Error())

				return
			}

			if payload.Status == "waiting" {
				err := saveToDatabaseUsers(dbConexao, payload)
				if err != nil {
					log.Printf("error no banco: %s", err)

					return
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
