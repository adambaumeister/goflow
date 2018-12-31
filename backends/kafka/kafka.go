package kafka

import (
	"crypto/tls"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/adambaumeister/goflow/fields"
	"os"
	"time"
)

type Kafka struct {
	server   string
	kconfig  *sarama.Config
	testMode bool
	tc       chan (string)

	producer sarama.AsyncProducer
}

func (b *Kafka) Configure(config map[string]string) {

	cd := map[string]string{
		"SSL":           "true",
		"SSL_VERIFY":    "false",
		"TEST_MODE":     "true",
		"SASL_USER":     "",
		"SASL_PASSWORD": "",
	}

	for k, v := range config {
		cd[k] = v
	}

	c := sarama.NewConfig()
	if val, ok := config["KAFKA_SERVER"]; ok {
		b.server = val
	} else if len(os.Getenv("KAFKA_SERVER")) > 0 {
		b.server = os.Getenv("KAFKA_SERVER")
	} else {
		panic("Invalid KAFKA Configuration. Missing KAFKA_SERVER.")
	}

	tls_config := tls.Config{}
	if cd["SSL"] == "true" {
		c.Net.TLS.Enable = true
		if cd["SSL_VERIFY"] != "true" {
			tls_config.InsecureSkipVerify = true
		}
		c.Net.TLS.Config = &tls_config
	}

	if cd["SASL_USER"] != "" {
		c.Net.SASL.Enable = true
		c.Net.SASL.User = cd["SASL_USER"]
		c.Net.SASL.Password = cd["SASL_PASSWORD"]
	}
	b.kconfig = c

	if cd["TEST_MODE"] == "true" {
		b.testMode = true
	}

}

//
func (b *Kafka) Init() {

	config := b.kconfig
	config.Producer.Return.Successes = true
	producer, err := sarama.NewAsyncProducer([]string{b.server}, config)
	if err != nil {
		panic(err)
	}

	b.producer = producer
	mt := make(map[uint16]fields.Value)
	b.Add(mt)
	// This is here to prevent Main from exiting preemptively when running Go Test.
	if b.testMode {
		time.Sleep(2 * time.Second)
	}
}

func (b *Kafka) Add(values map[uint16]fields.Value) {
	producer := b.producer

	// The idea is to try and read from producer.Errors() channel, if there's nothin' there, send the next message
	// Because add is called constantly, this will stop whenever an error is received.
	select {
	case err := <-producer.Errors():
		panic(fmt.Sprintf("Failed to produce message:: %v", err))
	default:
		producer.Input() <- &sarama.ProducerMessage{Topic: "test", Key: nil, Value: sarama.StringEncoder("Over SSL too!")}
	}
}
