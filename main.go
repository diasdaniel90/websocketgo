package main

import (
	"database/sql"
	"log"
	"runtime"
	"sync"

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
	go writePing(conn)
	log.Println("main", runtime.NumGoroutine())

	var wg sync.WaitGroup

	wg.Add(1)

	go controlMsg(&wg, conn, dbConexao, msgChan, errChan, msgStatusChan)

	wg.Wait()
}
