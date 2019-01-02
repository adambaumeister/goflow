package backends

import (
	"encoding/binary"
	"fmt"
	"github.com/adambaumeister/goflow/fields"
	"math/rand"
	"net"
)

//
// Backends are outbound interfaces for data
//
type Backend interface {
	Init()
	Status() string
	Configure(map[string]string)
	Add(map[uint16]fields.Value)
	Prune(string)
}

func GetTestFlow() map[uint16]fields.Value {
	testFlow := make(map[uint16]fields.Value)
	srcIP := fields.GetAddr(net.ParseIP("99.99.99.99"))
	dstIP := fields.GetAddr(net.ParseIP("99.99.99.99"))

	srcPortBytes := make([]byte, 2)
	dstPortBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(srcPortBytes, 19999)
	binary.BigEndian.PutUint16(dstPortBytes, 19999)
	srcPort := fields.GetInt(srcPortBytes)
	dstPort := fields.GetInt(dstPortBytes)

	protocol := fields.GetInt([]byte{6})
	srcPkts := fields.GetInt([]byte{254})
	srcBytes := fields.GetInt([]byte{254})

	testFlow[fields.IPV4_SRC_ADDR] = srcIP
	testFlow[fields.IPV4_DST_ADDR] = dstIP
	testFlow[fields.L4_SRC_PORT] = srcPort
	testFlow[fields.L4_DST_PORT] = dstPort
	testFlow[fields.PROTOCOL] = protocol
	testFlow[fields.IN_BYTES] = srcBytes
	testFlow[fields.IN_PKTS] = srcPkts
	v := fields.IntValue{Data: int(1546072176)}
	testFlow[fields.TIMESTAMP] = v
	return testFlow
}

func GetTestFlowRand(i int64) map[uint16]fields.Value {

	rand.Seed(i)

	ipstring := fmt.Sprintf("%v.%v.%v.%v", rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255))

	testFlow := make(map[uint16]fields.Value)
	srcIP := fields.GetAddr(net.ParseIP(ipstring))
	ipstring = fmt.Sprintf("%v.%v.%v.%v", rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255))
	dstIP := fields.GetAddr(net.ParseIP(ipstring))

	srcPortBytes := make([]byte, 2)
	dstPortBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(srcPortBytes, uint16(rand.Intn(65000)))
	binary.BigEndian.PutUint16(dstPortBytes, uint16(rand.Intn(65000)))
	srcPort := fields.GetInt(srcPortBytes)
	dstPort := fields.GetInt(dstPortBytes)

	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, rand.Uint32())
	protocol := fields.GetInt([]byte{6})
	srcPkts := fields.GetInt([]byte{254, 253})
	srcBytes := fields.GetInt([]byte{254, 253})

	testFlow[fields.IPV4_SRC_ADDR] = srcIP
	testFlow[fields.IPV4_DST_ADDR] = dstIP
	testFlow[fields.L4_SRC_PORT] = srcPort
	testFlow[fields.L4_DST_PORT] = dstPort
	testFlow[fields.PROTOCOL] = protocol
	testFlow[fields.IN_BYTES] = srcBytes
	testFlow[fields.IN_PKTS] = srcPkts
	v := fields.IntValue{Data: int(1546072176)}
	testFlow[fields.TIMESTAMP] = v
	return testFlow
}
