package kafka

import (
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// ProduceMessage envía un mensaje a Kafka
func ProduceMessage(topic string, key string, message string) error {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "kafka:9092",
	})
	if err != nil {
		return fmt.Errorf("error creando productor: %w", err)
	}
	defer p.Close()

	// Enviar mensaje
	err = p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            []byte(key),
		Value:          []byte(message),
	}, nil)

	if err != nil {
		return fmt.Errorf("error enviando mensaje: %w", err)
	}

	log.Printf("✅ Mensaje enviado a Kafka | Topic: %s | Key: %s | Message: %s", topic, key, message)

	return nil
}
