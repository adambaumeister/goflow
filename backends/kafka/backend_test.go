package kafka

import (
	"fmt"
	"github.com/adambaumeister/goflow/backends"
	"testing"
	"time"
)

const BENCH_MAX = 10

func DontTestBackend(t *testing.T) {
	fmt.Printf("Testing KAFKA Backend.\n")
	k := Kafka{}
	config := map[string]string{
		"TEST_MODE":   "true",
		"KAFKA_TOPIC": "test",
		"SASL_USER":   "admin",
	}

	k.Configure(config)
	k.Init()
	i := 0
	for i < BENCH_MAX {
		k.Add(backends.GetTestFlow())
		i++
	}
	time.Sleep(2 * time.Second)
}
