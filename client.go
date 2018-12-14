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
	FieldType2   int16
	FieldLength2 int16

	DataFlowSetID int16
	DataLength    int16
	DField1       uint8
	DField2       uint16
}

func Test() {
	conn, err := net.Dial("udp", "127.0.0.1:9999")
	if err != nil {
		fmt.Printf("Some error %v", err)
		return
	}
	tp := testPacket{
		// Header
		version:  9,
		count:    32,
		uptime:   1280,
		usecs:    122,
		sequence: 2000,
		id:       2300,
		// Template Headers
		flowSetId:  0,
		Length:     16,
		TemplateID: 2300,
		FieldCount: 2,
		// Template Fields
		FieldType1:   4,
		FieldLength1: 1,
		FieldType2:   7,
		FieldLength2: 2,
		// Data headers
		DataFlowSetID: 2300,
		DataLength:    7,
		// Data fields
		DField1: 6,
		DField2: 34477,
	}
	err = binary.Write(conn, binary.BigEndian, tp)
	if err != nil {
		fmt.Println("err:", err)
	}
}
