package main

import (
	"fmt"
	"net"
)

type packet struct {
	firstMessage  [1024]byte
	secondMessage [1024]byte
}

func main() {

	addr := net.UDPAddr{
		Port: 9999,
		IP:   net.ParseIP("127.0.0.1"),
	}
	conn, err := net.ListenUDP("udp", &addr)

	if err != nil {
		fmt.Printf("Some error %v\n", err)
		return
	}
	// Buffer creates an array of bytes
	buffer := make([]byte, 1024)
	// Read the number of bytes (1024) into a variable length slice of bytes, 'Buffer'
	count, _ := conn.Read(buffer)
	fmt.Println(string(buffer[:count]))
}
