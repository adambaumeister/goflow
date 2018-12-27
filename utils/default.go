package utils

import "github.com/adambaumeister/goflow/backends"

type Utility interface {
	SetBackends(map[string]backends.Backend)
	Run()
}
