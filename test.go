package main

import (
	"encoding/binary"
	"fmt"
	"net"
)

type packet struct {
	FirstMessage  [1024]byte
	SecondMessage [1024]byte
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
	// 'Count' refers to the number of bytes received in the slice
	// Below we decode that amount (buffer_slice[:number_of_bytes]) as a string
	// fmt.Println(string(buffer[:count]))

	p := packet{}
	binary.Read(*buffer)

}