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
	FlowSetID uint16
	Length    uint16
	ID        uint16
}

func main() {

	// Provides parsing for Netflow V9 Records
	// https://www.ietf.org/rfc/rfc3954.txt

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
	// Buffer creates an array of bytes
	// We want to read the entire datagram in as UDP is type SOCK_DGRAM and "Read" can't be called more than once
	packet := make([]byte, 1500)
	// Read the max number of bytes in a datagram(1500) into a variable length slice of bytes, 'Buffer'
	conn.Read(packet)
	id := netflowPacketFlowsetId{}
	p.Version = binary.BigEndian.Uint16(packet[:2])
	id.FlowSetID = binary.BigEndian.Uint16(packet[20:22])

	switch id.FlowSetID {
	// Template flowset
	case 0:
		template := netflowPacketTemplate{}
		template.Length = binary.BigEndian.Uint16(packet[22:24])
		templateSlice := packet[22:template.Length]
		template.ID = binary.BigEndian.Uint16(templateSlice[2:4])
		fmt.Printf("version: %v Template %v", p.Version, template.ID)
	}
}
