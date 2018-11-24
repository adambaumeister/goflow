package frontends

import "github.com/adamb/goflow/backends"

//
// Frontends are inbound interfaces for data
//
type Frontend interface {
	Start()
	Configure(config map[string]string, backend backends.Backend)
}
