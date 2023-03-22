package main

import "net"

const address = "127.0.0.1:20001"

func sendUDPMessage(message string) error {
	serverAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return err
	}
	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Write([]byte(message))
	if err != nil {
		return err
	}

	return nil
}
