package main

import (
	"encoding/binary"
	"fmt"
	"net"
)

// CONSTANTS
const IN_BYTES = 1
const IN_PKTS = 2
const FLOWS = 3
const PROTOCOL = 4
const TOS = 5
const TCP_FLAGS = 6
const L4_SRC_PORT = 7

var FUNCTIONMAP = map[uint8]func([]byte) interface{}{
	IN_BYTES: GetInt,
}

// GENERICS
type netflow struct {
	Templates map[uint16]netflowPacketTemplate
}
type netflowPacket struct {
	Header    netflowPacketHeader
	Length    int
	Templates map[uint16]netflowPacketTemplate
	Data      []netflowDataFlowset
}
type netflowPacketHeader struct {
	Version  uint16
	Count    int16
	Uptime   int32
	Sequence int32
	Id       int32
}
type netflowPacketFlowset struct {
	FlowSetID uint16
	Length    uint16
}

// TEMPLATE STRUCTS
type netflowPacketTemplate struct {
	FlowSetID  uint16
	Length     uint16
	ID         uint16
	FieldCount uint16
	Fields     []templateField
}
type templateField struct {
	FieldType uint16
	Length    uint16
}

// DATA STRUCTS
type netflowDataFlowset struct {
	FlowSetID uint16
	Length    uint16
	Records   []flowRecord
}
type flowRecord struct {
	Values []interface{}
}

func GetInt(p []byte) interface{} {
	switch {
	case len(p) > 1:
		return binary.BigEndian.Uint16(p)
	}
}

/*
ParseData

Takes a slice of a data flowset and retreives all the flow records
Requires
	p []byte : Data Flowset slice
*/
func parseData(n netflowPacket, p []byte) netflowDataFlowset {
	nfd := netflowDataFlowset{
		FlowSetID: binary.BigEndian.Uint16(p[:2]),
		Length:    binary.BigEndian.Uint16(p[2:4]),
	}

	t := n.Templates[nfd.FlowSetID]
	start := uint16(4)
	// Read each Field in order from the flowset until the length is exceeded
	for start < nfd.Length {
		fr := flowRecord{}
		for _, f := range t.Fields {
			value := p[start : start+f.Length]
			switch f.FieldType {
			case 4:
				fmt.Printf("Field 4 value: %v\n", uint8(value[0]))
			case 7:
				fmt.Printf("Field  value: %v\n", binary.BigEndian.Uint16(value))
			}

			fr.Values = append(fr.Values, value)
			start = start + f.Length
		}
		nfd.Records = append(nfd.Records, fr)
	}

	return nfd
}

/*
ParseTemplate

Slices a flow template out of an overall packet
Requires
	p []byte : Full packet bytes
Returns
	netFlowPacketTemplate: Struct of template
*/
func parseTemplate(templateSlice []byte) netflowPacketTemplate {

	template := netflowPacketTemplate{
		Fields: make([]templateField, 0),
	}
	template.ID = binary.BigEndian.Uint16(templateSlice[4:6])

	// Get the number of Fields
	template.FieldCount = binary.BigEndian.Uint16(templateSlice[6:8])
	// Start at the first fields and work through
	fmt.Printf("L; %v, I: %v, FC: %v\n\n", template.Length, template.ID, template.FieldCount)
	fieldStart := 8
	var read = uint16(0)
	for read < template.FieldCount {
		fieldTypeEnd := fieldStart + 2
		fieldType := binary.BigEndian.Uint16(templateSlice[fieldStart:fieldTypeEnd])
		fieldLengthEnd := fieldTypeEnd + 2
		fieldLength := binary.BigEndian.Uint16(templateSlice[fieldTypeEnd:fieldLengthEnd])

		// Create template FIELD struct
		field := templateField{
			FieldType: fieldType,
			Length:    fieldLength,
		}
		// Template fields are IN ORDER
		// Order determines records in data flowset
		template.Fields = append(template.Fields, field)

		read++
		fieldStart = fieldLengthEnd
	}

	for _, template := range template.Fields {
		fmt.Printf("Type read: %v, Length: %v\n", template.FieldType, template.Length)
	}
	return template
}

/*
Route
Takes an entire packet slice, and routes each flowset to the correct handler

Requires
	netflowPacket : netflowpacket struct
	[]byte		  : Packet bytes
	uint16		  : Byte index to start at (skip the headers, etc)
*/
func Route(nfp netflowPacket, p []byte, start uint16) {
	id := uint16(0)
	l := uint16(0)

	//fmt.Printf("End of slice: %v", start+id.Length)
	for int(start) < nfp.Length {
		id = binary.BigEndian.Uint16(p[start : start+2])
		l = binary.BigEndian.Uint16(p[start+2 : start+4])
		// Slice the next flowset out
		s := p[start : start+l]
		// Flowset ID is the switch we use to determine what sort of flowset follors
		switch {
		// Template flowset
		case id == uint16(0):
			t := parseTemplate(s)
			nfp.Templates[t.ID] = t
		// Data flowset
		case id > uint16(255):
			d := parseData(nfp, s)
			nfp.Data = append(nfp.Data, d)
		}
		start = start + l
	}
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

	nfpacket := netflowPacket{
		Templates: make(map[uint16]netflowPacketTemplate),
	}

	p := netflowPacketHeader{}
	// Buffer creates an array of bytes
	// We want to read the entire datagram in as UDP is type SOCK_DGRAM and "Read" can't be called more than once
	packet := make([]byte, 1500)
	// Read the max number of bytes in a datagram(1500) into a variable length slice of bytes, 'Buffer'
	// Also set the total number of bytes read so we can check it later
	nfpacket.Length, _ = conn.Read(packet)
	fmt.Printf("Total packet length: %v\n", nfpacket.Length)
	p.Version = binary.BigEndian.Uint16(packet[:2])
	nfpacket.Header = p

	Route(nfpacket, packet, uint16(20))
}
