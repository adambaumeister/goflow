package fields

import (
	"fmt"
	"net"
)

// Retrieve an addr value from a field
func GetAddr(p []byte) Value {
	var a AddrValue
	var ip net.IP
	ip = p
	a.Data = ip
	return a
}

//
// Address Values
//
type AddrValue struct {
	Data net.IP
}

func (a AddrValue) ToString() string {
	return fmt.Sprintf("%v", a.Data.String())
}
