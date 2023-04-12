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
	msgSignalChan := make(chan MsgSignal)

	go listenUDP(msgSignalChan)

	dbConexao, err := sql.Open("mysql", EnvsDatabase())
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
	msgStatusChan := make(chan MsgStatus)

	go BetStatus(msgStatusChan, msgSignalChan)

	go readMessages(conn, msgChan, errChan)
	go writePing(conn)
	log.Println("main", runtime.NumGoroutine())

	var wg sync.WaitGroup

	wg.Add(1)

	go ControlMsg(&wg, conn, dbConexao, msgChan, errChan, msgStatusChan)

	wg.Wait()
}

func BetStatus(msgStatusChan <-chan MsgStatus, msgSignalChan <-chan MsgSignal) {
	log.Println("###########11##################")

	mensagens := []MsgSignal{}
	// var valido string
	for {
		select {
		case msg, ok := <-msgStatusChan:
			if !ok {
				log.Println("Canal msgStatusChan fechado")

				return
			}

			log.Println("msgStatusChan", msg)

			// if valido == msg.BetStatus {
			// }
			// if msg.BetStatus == waiting {
			// 	log.Println("------+++++++++#######vai esperar 6 segundos ")
			// 	time.Sleep(6 * time.Second)
			// }

			time.AfterFunc(tempoEspera*time.Second, func() {
				ControBet(mensagens)
				// log.Println("+++++Recebeu os sinais na msgStatusChan ", mensagens)
			})

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

func ControBet(mensagens []MsgSignal) {
	log.Println("Executando a função após 4 segundos...", mensagens)
}

func ControlMsg(wg *sync.WaitGroup, conn io.Closer, dbConexao *sql.DB, msgChan chan []byte,
	errChan chan error, msgStatusChan chan MsgStatus,
) {
	defer wg.Done()

	var lastMsg LastMsg

	for {
		select {
		case msg := <-msgChan:
			payload, err := decodePayload(msg[2:])
			if err != nil {
				log.Printf("Erro ao decodificar mensagem: %s", err)

				return
			}

			Status, err := filterMessage(dbConexao, payload, &lastMsg)
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
