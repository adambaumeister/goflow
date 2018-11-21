package fields

import (
	"encoding/binary"
	"fmt"
)

//
// Integer Values
//
type IntValue struct {
	Data int
	Type uint16
}

func (i IntValue) SetType(t uint16) {
	i.Type = t
}
func (i IntValue) ToString() string {
	return fmt.Sprintf("%v", i.Data)
}

// Retrieve integer values from a field
func GetInt(p []byte) Value {
	var i IntValue
	switch {
	case len(p) > 2:
		i.Data = int(binary.BigEndian.Uint32(p))
		return i
	case len(p) > 1:
		i.Data = int(binary.BigEndian.Uint16(p))
		return i
	default:
		i.Data = int(uint8(p[0]))
		return i
	}
}
