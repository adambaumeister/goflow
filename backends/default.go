package backends

import "github.com/adambaumeister/goflow/fields"

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
