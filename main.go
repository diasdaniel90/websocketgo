package main

import (
	"database/sql"
	"log"
	"runtime"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

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

	msgChanWebsocket := make(chan []byte)
	errChan := make(chan error)
	msgStatusChan := make(chan msgStatusStruct)

	go controlBet(dbConexao, msgStatusChan, msgSignalChan)

	go readMessages(conn, msgChanWebsocket, errChan)
	go writePing(conn)
	log.Println("main", runtime.NumGoroutine())

	var wg sync.WaitGroup

	wg.Add(1)

	go controlMsg(&wg, conn, dbConexao, msgChanWebsocket, errChan, msgStatusChan)

	wg.Wait()
}
