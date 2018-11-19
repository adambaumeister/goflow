package main

import (
	"encoding/binary"
	"fmt"
	"net"
)

type netflow struct {
	Templates map[uint16]netflowPacketTemplate
}

type netflowPacket struct {
	Header    netflowPacketHeader
	Length    int
	Templates map[uint16]netflowPacketTemplate
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

type netflowPacketTemplate struct {
	FlowSetID  uint16
	Length     uint16
	ID         uint16
	FieldCount uint16
	Fields     map[uint16]templateField
}

type templateField struct {
	FieldType uint16
	Length    uint16
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
		Fields: make(map[uint16]templateField),
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
		// Assign to fields, keyed by the type, to this template record
		template.Fields[fieldType] = field

		read++
		fieldStart = fieldLengthEnd
	}

	for t, template := range template.Fields {
		fmt.Printf("Type read: %v, Length: %v\n", t, template.Length)
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
	for int(start+l) <= nfp.Length {
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
			fmt.Printf("Yes")
		}
		start = start + l
		fmt.Printf("New length: %v", start)
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
