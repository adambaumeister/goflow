package frontends

import (
	"github.com/adamb/goflow/backends"
)

//
// Frontends are inbound interfaces for data
//
type Frontend interface {
	Start(backend backends.Backend)
	Configure(config map[string]string)
}
