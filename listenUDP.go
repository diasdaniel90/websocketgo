package main

import (
	"encoding/json"
	"log"
	"net"
)

const (
	host = "127.0.0.1"
	port = 1234
	size = 1024
)

func listenUDP(msgSignalChan chan MsgSignal) {
	ipv4 := net.ParseIP(host)

	conn, err := net.ListenUDP("udp", &net.UDPAddr{IP: ipv4, Port: port, Zone: ""})
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	log.Println("Listening on", conn.LocalAddr().String())

	buf := make([]byte, size)

	for {
		nBytes, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Println("Error:", err)

			continue
		}

		// log.Printf("Received %d bytes from %s: %s\n", nBytes, addr.String(), string(buf[:nBytes]))

		var msgSignal MsgSignal

		err = json.Unmarshal((buf[:nBytes]), &msgSignal)
		if err != nil {
			log.Printf("Erro ao ler buffer UDP: %s", err.Error())
			panic(err.Error())
		}

		// log.Println(msgSignal)
		msgSignalChan <- msgSignal
	}
}
