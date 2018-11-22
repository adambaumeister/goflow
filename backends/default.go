package backends

import "github.com/adamb/goflow/fields"

//
// Backends are outbound interfaces for data
//
type Backend interface {
	Init()
	Test()
	Add(map[uint16]fields.Value)
}
