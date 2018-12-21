package fields

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"
	"net"
)

//
// Address Values
//
// v4
type AddrValue struct {
	Data net.IP
	Type uint16
	Int  uint32
}

func (i AddrValue) ToInt() int {
	return int(i.Int)
}

func (i AddrValue) SetType(t uint16) {
	i.Type = t
}
func (a AddrValue) ToString() string {
	return fmt.Sprintf("%v", a.Data.String())
}

//v6
type Addr6Value struct {
	Data string
	Type uint16
	Int  *big.Int
}

func (i Addr6Value) ToInt() int {
	// V6 addresses don't fit in a 64-bit UINT so this function is uncallable
	return 0
}

func (i Addr6Value) SetType(t uint16) {
	i.Type = t
}
func (a Addr6Value) ToString() string {
	return fmt.Sprintf("%v", a.Data)
}

// Retrieve an addr value from a field
func GetAddr(p []byte) Value {
	var a AddrValue
	var ip net.IP
	ip = p
	a.Data = ip
	a.Int = binary.BigEndian.Uint32(p)
	return a
}

func GetAddr6(p []byte) Value {
	var a Addr6Value
	bi := big.NewInt(0)
	bi.SetBytes(p)

	a.Int = bi
	a.Data = hex.EncodeToString(p)
	return a
}
