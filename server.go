package main

import (
	"encoding/binary"
	"fmt"
	"net"
)

type netflowPacketHeader struct {
	Version   uint16
	Count     int16
	Uptime    int32
	Sequence  int32
	Id        int32
	FlowSetId netflowPacketFlowsetId
	Length    netflowPacketTemplate
}

type netflowPacketFlowsetId struct {
	FlowSetID uint16
}

type netflowPacketTemplate struct {
	Length int16
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
	//f := netflowPacketFlowsetId{}
	// Buffer creates an array of bytes
	// We want to read the entire datagram in as UDP is type SOCK_DGRAM and "Read" can't be called more than once
	packet := make([]byte, 1500)
	// Read the number of bytes (1500) into a variable length slice of bytes, 'Buffer'
	count, _ := conn.Read(packet)
	id := netflowPacketFlowsetId{}
	p.Version = binary.BigEndian.Uint16(packet[:2])
	id.FlowSetID = binary.BigEndian.Uint16(packet[20:22])
	fmt.Printf("version: %v FlowSetID %v(read: %v)", p.Version, id.FlowSetID, count)

	// 'Count' refers to the number of bytes received in the slice
	// Below we decode that amount (buffer_slice[:number_of_bytes]) as a string
	// fmt.Println(string(buffer[:count]))

}
