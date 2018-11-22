package fields

import (
	"encoding/binary"
	"fmt"
	"net"
)

//
// Address Values
//
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

// Retrieve an addr value from a field
func GetAddr(p []byte) Value {
	var a AddrValue
	var ip net.IP
	ip = p
	a.Data = ip
	a.Int = binary.BigEndian.Uint32(p)
	return a
}
