package fields

import (
	"fmt"
	"net"
)

//
// Address Values
//
type AddrValue struct {
	Data net.IP
	Type uint16
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
	return a
}
