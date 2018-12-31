package kafka

import "testing"

func TestBackend(t *testing.T) {
	k := Kafka{}
	config := make(map[string]string)
	k.Configure(config)
	k.Init()
}
