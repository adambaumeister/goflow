package frontends

import (
	"encoding/binary"
	"fmt"
	"github.com/adambaumeister/goflow/backends"
	"github.com/adambaumeister/goflow/fields"
	"net"
	"os"
	"strconv"
	"strings"
)

// CONSTANTS
// Actual fields...
const IN_BYTES = fields.IN_BYTES
const IN_PKTS = fields.IN_PKTS
const FLOWS = fields.FLOWS
const PROTOCOL = fields.PROTOCOL
const TOS = fields.TOS
const TCP_FLAGS = fields.TCP_FLAGS
const L4_SRC_PORT = fields.L4_SRC_PORT
const IPV4_SRC_ADDR = fields.IPV4_SRC_ADDR
const SRC_MASK = fields.SRC_MASK
const L4_DST_PORT = fields.L4_DST_PORT
const IPV4_DST_ADDR = fields.IPV4_DST_ADDR
const IPV4_NEXT_HOP = fields.IPV4_NEXT_HOP
const OUT_BYTES = fields.OUT_BYTES
const OUT_PKTS = fields.OUT_PKTS
const LAST_SWITCHED = fields.LAST_SWITCHED

// Extension fields
const _TIMESTAMP = fields.TIMESTAMP

var FUNCTIONMAP = map[uint16]func([]byte) fields.Value{
	IN_BYTES:      fields.GetInt,
	IN_PKTS:       fields.GetInt,
	PROTOCOL:      fields.GetInt,
	L4_SRC_PORT:   fields.GetInt,
	IPV4_SRC_ADDR: fields.GetAddr,
	IPV4_DST_ADDR: fields.GetAddr,
	OUT_BYTES:     fields.GetInt,
	OUT_PKTS:      fields.GetInt,
	L4_DST_PORT:   fields.GetInt,
	LAST_SWITCHED: fields.GetInt,
}

//
// GENERICS
//
// Netflow listener and main object
type Netflow struct {
	Templates map[uint16]netflowPacketTemplate
	BindAddr  net.IP
	BindPort  int

	backend backends.Backend
}
type netflowPacket struct {
	Header    netflowPacketHeader
	Length    int
	Templates map[uint16]netflowPacketTemplate
	Data      []netflowDataFlowset
}
type netflowPacketHeader struct {
	Version  uint16
	Count    uint16
	Uptime   uint32
	Usecs    uint32
	Sequence uint32
	Id       uint32
}
type netflowPacketFlowset struct {
	FlowSetID uint16
	Length    uint16
}

// TEMPLATE STRUCTS
type netflowPacketTemplate struct {
	FlowSetID   uint16
	Length      uint16
	ID          uint16
	FieldCount  uint16
	Fields      []templateField
	FieldLength uint16
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
	Values    []fields.Value
	ValuesMap map[uint16]fields.Value
}

func (r *flowRecord) calcTime(s uint32, u uint32) uint32 {
	/*
		Calculate the timestamp of a record end time using the following:
		Usecs - SysUptime + FlowEndTime

	*/
	var ts uint32

	if flowendSecs, ok := r.ValuesMap[LAST_SWITCHED]; ok {
		ts = u - (s + uint32(flowendSecs.ToInt()))
		//ts = u
		v := fields.IntValue{Data: int(ts)}
		r.ValuesMap[_TIMESTAMP] = v
	}
	return ts
}
func (r flowRecord) toString() string {
	var sl []string
	for _, v := range r.Values {
		sl = append(sl, v.ToString())
	}
	return strings.Join(sl, " : ") + "\n"
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

	// Return no flow records if it's empty
	if _, ok := n.Templates[nfd.FlowSetID]; !ok {
		return nfd
	}
	t := n.Templates[nfd.FlowSetID]

	start := uint16(4)
	// Read each Field in order from the flowset until the length is exceeded
	for start < nfd.Length {
		// Check the number of fields don't overrun the size of this flowset
		// if so, remainder must be padding
		if t.FieldLength <= (nfd.Length - start) {
			fr := flowRecord{ValuesMap: make(map[uint16]fields.Value)}
			for _, f := range t.Fields {
				valueSlice := p[start : start+f.Length]
				if function, ok := FUNCTIONMAP[f.FieldType]; ok {
					value := function(valueSlice)
					value.SetType(f.FieldType)
					fr.Values = append(fr.Values, value)
					fr.ValuesMap[f.FieldType] = value
				}

				start = start + f.Length
			}
			nfd.Records = append(nfd.Records, fr)
		} else {
			start = start + (nfd.Length - t.FieldLength)
		}
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
		template.FieldLength = template.FieldLength + fieldLength
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
func Route(nfp netflowPacket, p []byte, start uint16) netflowPacket {
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
	return nfp
}

func (n *Netflow) Configure(config map[string]string, b backends.Backend) {
	/*
		Configure
		Configures this object
		Requires
			config : K/V map of configuration options
	*/
	var test int
	n.BindAddr = net.ParseIP(config["bindaddr"])
	test, err := strconv.Atoi(config["bindport"])
	n.BindPort = test
	if err != nil {
		panic(err.Error())
	}
	n.backend = b
}

func (nf Netflow) Start() {
	// Provides parsing for Netflow V9 Records
	// https://www.ietf.org/rfc/rfc3954.txt
	b := nf.backend
	b.Init()
	addr := net.UDPAddr{
		Port: nf.BindPort,
		IP:   nf.BindAddr,
	}
	conn, err := net.ListenUDP("udp", &addr)

	if err != nil {
		fmt.Printf("Some error %v\n", err)
		return
	}
	fmt.Printf("Listen on Addr: %v, Port: %v", nf.BindAddr, nf.BindPort)
	nf.Templates = make(map[uint16]netflowPacketTemplate)
	// Listen to incoming flows
	for {
		nfpacket := netflowPacket{
			Templates: nf.Templates,
		}

		p := netflowPacketHeader{}
		// Buffer creates an array of bytes
		// We want to read the entire datagram in as UDP is type SOCK_DGRAM and "Read" can't be called more than once
		packet := make([]byte, 1500)
		// Read the max number of bytes in a datagram(1500) into a variable length slice of bytes, 'Buffer'
		// Also set the total number of bytes read so we can check it later
		nfpacket.Length, _ = conn.Read(packet)

		p.Version = binary.BigEndian.Uint16(packet[:2])
		p.Uptime = binary.BigEndian.Uint32(packet[4:8])
		p.Usecs = binary.BigEndian.Uint32(packet[8:12])
		switch p.Version {
		case 5:
			fmt.Printf("Wrong Netflow version, only v9 supported.")
			os.Exit(1)
		}
		nfpacket.Header = p

		nfpacket = Route(nfpacket, packet, uint16(20))
		nf.Templates = nfpacket.Templates
		for _, dfs := range nfpacket.Data {
			for _, record := range dfs.Records {
				record.calcTime(p.Uptime, p.Usecs)
				b.Add(record.ValuesMap)
			}
		}
	}
}
