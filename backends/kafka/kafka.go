package kafka

import (
	"crypto/tls"
	"fmt"
	"github.com/Shopify/sarama"
	"log"
	"os"
	"os/signal"
	"time"
)

type Kafka struct {
	server  string
	kconfig *sarama.Config

	producer *sarama.AsyncProducer
}

func (b *Kafka) Configure(config map[string]string) {

	cd := map[string]string{
		"SSL":        "true",
		"SSL_VERIFY": "false",
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
			tls_config.InsecureSkipVerify = false
		}
		c.Net.TLS.Config = &tls_config
	}
	b.kconfig = c
}

//
func (b *Kafka) Init() {

	config := sarama.NewConfig()
	config.Net.TLS.Enable = true
	config.Net.TLS.Config = &tls.Config{InsecureSkipVerify: true}
	config.Producer.Return.Successes = true
	producer, err := sarama.NewAsyncProducer([]string{b.server}, config)
	if err != nil {
		panic(err)
	}

	// Trap SIGINT to trigger a shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	select {
	case producer.Input() <- &sarama.ProducerMessage{Topic: "test", Key: nil, Value: sarama.StringEncoder("Over SSL too!")}:
		fmt.Printf("Sent")
	case err := <-producer.Errors():
		log.Println("Failed to produce message", err)
		break
	case <-signals:
		break
	}
	time.Sleep(5 * time.Second)
}
