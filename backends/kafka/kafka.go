package kafka

import (
	"github.com/Shopify/sarama"
	"log"
	"os"
	"os/signal"
)

type Kafka struct {
}

//
func (b *Kafka) Init() {

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	producer, err := sarama.NewAsyncProducer([]string{"ec2-13-210-246-197.ap-southeast-2.compute.amazonaws.com:9092"}, config)
	if err != nil {
		panic(err)
	}

	var enqueued, errors int

	// Trap SIGINT to trigger a shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
ProducerLoop:
	for {
		select {
		case producer.Input() <- &sarama.ProducerMessage{Topic: "test", Key: nil, Value: sarama.StringEncoder("testing 123")}:
			enqueued++
		case err := <-producer.Errors():
			log.Println("Failed to produce message", err)
			errors++
			break ProducerLoop
		case <-signals:
			break ProducerLoop
		}
	}

	log.Printf("Enqueued: %d; errors: %d\n", enqueued, errors)
}
