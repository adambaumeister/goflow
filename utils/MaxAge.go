package utils

import (
	"fmt"
	"github.com/adambaumeister/goflow/backends"
	"time"
)

/*
MaxAge
Using exposed "prune" methods from backends and associated GlobalConfig, deletes flow records.
*/
type MaxAge struct {
	// String because it's going into a query string anyway
	MaxAgeDays string
	backends   map[string]backends.Backend
}

// Set the backends to operate on
func (m *MaxAge) SetBackends(b map[string]backends.Backend) {
	m.backends = b
}

func (m *MaxAge) Run() {
	for {
		for _, be := range m.backends {
			be.Prune(m.MaxAgeDays)
		}
		fmt.Printf("Pruning configured backends...")
		// This is really gross, we should have some sort of internal crontab-like thing for utilities
		time.Sleep(86400 * time.Second)
	}
}
