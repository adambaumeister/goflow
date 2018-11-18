package main

import (
	"encoding/binary"
	"fmt"
	"net"
)

type testPacket struct {
	version    int16
	count      int16
	uptime     int32
	usecs      int32
	sequence   int32
	id         int32
	flowSetId  int16
	Length     int16
	TemplateID int16
	FieldCount int16

	FieldType1   int16
	FieldLength1 int16
}

func main() {
	conn, err := net.Dial("udp", "127.0.0.1:9999")
	if err != nil {
		fmt.Printf("Some error %v", err)
		return
	}
	tp := testPacket{
		version:    9,
		count:      32,
		uptime:     1280,
		usecs:      122,
		sequence:   2000,
		id:         2300,
		flowSetId:  0,
		Length:     77,
		TemplateID: 11,
		FieldCount: 1,

		FieldType1:   1,
		FieldLength1: 32,
	}
	err = binary.Write(conn, binary.BigEndian, tp)
	if err != nil {
		fmt.Println("err:", err)
	}
}
