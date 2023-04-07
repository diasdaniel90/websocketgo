package main

import (
	"log"
	"net"
)

const (
	Host = "127.0.0.1"
	Port = 1234
	Size = 1024
)

func listenUDP() {
	ipv4 := net.ParseIP(Host)

	conn, err := net.ListenUDP("udp", &net.UDPAddr{IP: ipv4, Port: Port, Zone: ""})
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	log.Println("Listening on", conn.LocalAddr().String())

	buf := make([]byte, Size)

	for {
		nBytes, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Println("Error:", err)

			continue
		}

		log.Printf("Received %d bytes from %s: %s\n", nBytes, addr.String(), string(buf[:nBytes]))
	}
}