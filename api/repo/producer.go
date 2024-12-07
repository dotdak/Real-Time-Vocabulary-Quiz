package repo

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM/sarama"
)

type Producer struct {
	topic string
	sarama.AsyncProducer
}

func (p *Producer) Send(data interface{}) {
	payload, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("can not marshal payload %v", err)
	}
	// Produce a message to a Kafka topic
	message := &sarama.ProducerMessage{
		Topic: p.topic, // Replace with your topic
		Value: sarama.ByteEncoder(payload),
	}

	// Send the message
	p.Input() <- message

	// Handle success and error messages
	select {
	case msg := <-p.Successes():
		response, _ := msg.Value.Encode()
		fmt.Printf("Message sent successfully: %v\n", string(response))
	case err := <-p.Errors():
		log.Printf("Error sending message: %v\n", err)
	}
}

func (p *Producer) Close() {
	p.AsyncClose()
}

func NewProducer(topic string) *Producer {
	// Define Kafka broker addresses
	brokers := []string{"localhost:9092"} // Replace with your broker address

	// Create a new Sarama configuration
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	// Create a new Kafka producer
	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		log.Fatalf("Error creating Kafka producer: %v", err)
	}

	return &Producer{
		AsyncProducer: producer,
		topic:         topic,
	}
}
