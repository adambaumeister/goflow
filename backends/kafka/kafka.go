package kafka

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/adambaumeister/goflow/fields"
	"os"
	"testing"
)

type Kafka struct {
	server   string
	kconfig  *sarama.Config
	testMode bool
	tc       chan (string)
	topic    string

	producer sarama.AsyncProducer
}

type JsonFLow struct {
	Src_ip     string
	Dst_ip     string
	Src_port   int
	Dst_port   int
	Protocol   int
	In_bytes   int
	In_packets int
	Src_ip6    string
	Dst_ip6    string
}

func (j *JsonFLow) route(values map[uint16]fields.Value) {
	// There's probably a nicer way of doing this.
	for f, v := range values {
		switch f {
		case fields.IPV4_SRC_ADDR:
			j.Src_ip = v.ToString()
		case fields.IPV4_DST_ADDR:
			j.Dst_ip = v.ToString()
		case fields.L4_SRC_PORT:
			j.Src_port = v.ToInt()
		case fields.L4_DST_PORT:
			j.Dst_port = v.ToInt()
		case fields.PROTOCOL:
			j.Protocol = v.ToInt()
		case fields.IN_BYTES:
			j.In_bytes = v.ToInt()
		case fields.IN_PKTS:
			j.In_packets = v.ToInt()
		case fields.IPV6_SRC_ADDR:
			j.Src_ip6 = v.ToString()
		case fields.IPV6_DST_ADDR:
			j.Dst_ip6 = v.ToString()
		}
	}
}

func (b *Kafka) Prune(string) {
}

func (b *Kafka) BenchmarkBackend(t testing.B) {
}

func (b *Kafka) Status() string {
	b.Init()
	return "Kafka connection looks ok.\n"
}

func (b *Kafka) Configure(config map[string]string) {

	// Config required - no defaults
	cr := []string{
		"KAFKA_SERVER",
		"KAFKA_TOPIC",
	}
	// Config defauls - Optional arguments - they have defaults
	cd := map[string]string{
		"SSL":           "true",
		"SSL_VERIFY":    "false",
		"TEST_MODE":     "true",
		"SASL_USER":     "",
		"SASL_PASSWORD": "",
	}

	// Overwrite the defaults with the real values
	for k, v := range config {
		cd[k] = v
	}

	//  Overwrite any values with those set in the environment, if existing
	for k, _ := range cd {
		if len(os.Getenv(k)) > 0 {
			cd[k] = os.Getenv(k)
		}
	}

	c := sarama.NewConfig()
	for _, v := range cr {
		if _, ok := config[v]; !ok {
			if len(os.Getenv(v)) > 0 {
				config[v] = os.Getenv(v)
			} else {
				panic(fmt.Sprintf("Invalid Kafka Configuration. Missing %v", v))
			}
		}
	}
	b.server = config["KAFKA_SERVER"]
	b.topic = config["KAFKA_TOPIC"]

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
	//config.Producer.Return.Successes = true
	producer, err := sarama.NewAsyncProducer([]string{b.server}, config)
	if err != nil {
		panic(err)
	}

	b.producer = producer
}

func (b *Kafka) Add(values map[uint16]fields.Value) {
	producer := b.producer

	jf := JsonFLow{}
	jf.route(values)

	s, err := json.Marshal(jf)
	if err != nil {
		panic(fmt.Sprintf("Failed to marshal JSON Flow: %v", err))
	}

	// The idea is to try and read from producer.Errors() channel, if there's nothin' there, send the next message
	// Because add is called constantly, this will stop whenever an error is received.
	select {
	case err := <-producer.Errors():
		panic(fmt.Sprintf("Failed to produce message: %v", err))
	default:
		producer.Input() <- &sarama.ProducerMessage{Topic: "test", Key: nil, Value: sarama.ByteEncoder(s)}
	}
	// This is here to prevent Main from exiting preemptively when running Go Test.
	if b.testMode {
		//time.Sleep(2 * time.Second)
	}
}
