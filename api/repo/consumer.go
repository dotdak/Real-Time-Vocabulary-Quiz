package repo

import (
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
)

func NewConsumer[T any](topic string, queue chan<- T) {
	brokers := []string{"localhost:9092"}
	// Create a new Sarama configuration
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	// Create a new Kafka consumer
	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Fatalf("Error creating Kafka consumer: %v", err)
	}

	partitionList, err := consumer.Partitions(topic)
	if err != nil {
		log.Fatal("Failed to get partitions:", err)
	}

	for _, partition := range partitionList {
		pc, err := consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
		if err != nil {
			log.Fatal("Failed to start consumer for partition", partition, ":", err)
		}

		go func(pc sarama.PartitionConsumer) {
			defer pc.Close()
			for msg := range pc.Messages() {
				log.Printf("Received message in topic %s: %s\n", topic, string(msg.Value))
				var data T
				if err := json.Unmarshal(msg.Value, &data); err != nil {
					log.Fatalf("cannot unmarshal message: %v", err)
				}

				queue <- data
			}
		}(pc)
	}

	// Subscribe to a Kafka topic
	// partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	// if err != nil {
	// 	log.Fatalf("Error consuming Kafka partition: %v", err)
	// }

	// go func() {
	// 	defer consumer.Close()
	// 	defer partitionConsumer.Close()
	// 	// Consume messages from Kafka
	// 	for msg := range partitionConsumer.Messages() {
	// 		log.Printf("Received message: %s\n", string(msg.Value))
	// 		var data T
	// 		if err := json.Unmarshal(msg.Value, &data); err != nil {
	// 			log.Fatalf("cannot unmarshal message: %v", err)
	// 		}

	// 		queue <- data
	// 	}
	// }()
}
