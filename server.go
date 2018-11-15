package main

import (
	"encoding/binary"
	"fmt"
	"net"
)

type netflowPacketHeader struct {
	Version  int16
	Count    int16
	Uptime   int32
	Sequence int32
	Id       int32
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

	p := netflowPacketHeader{}
	err = binary.Read(conn, binary.BigEndian, &p)
	if err != nil {
		fmt.Printf("Some error %v\n", err)
		return
	}
	fmt.Printf("Int: %v %v\n", p.Version, p.Count)
	// Buffer creates an array of bytes
	//buffer := make([]byte, 1024)

	// Read the number of bytes (1024) into a variable length slice of bytes, 'Buffer'
	//count, _ := conn.Read(buffer)

	// 'Count' refers to the number of bytes received in the slice
	// Below we decode that amount (buffer_slice[:number_of_bytes]) as a string
	// fmt.Println(string(buffer[:count]))

}
