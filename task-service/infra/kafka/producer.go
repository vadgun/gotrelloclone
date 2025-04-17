package kafka

import (
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// ProduceMessage envía un mensaje o publica el evento en Kafka
func ProduceMessage(userID, message, topic, key string) error {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "kafka:9092",
	})
	if err != nil {
		return fmt.Errorf("❌ Error al crear productor de Kafka: %v", err)
	}
	defer p.Close()

	// Enviar el mensaje o publicamos el evento en Kafka
	err = p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            []byte(key),
		Value:          []byte(message),
	}, nil)

	if err != nil {
		return fmt.Errorf("❌ Error enviando mensajea Kafka: %w", err)
	}

	log.Printf("✅ Mensaje enviado a Kafka | Topic: %s | Key: %s | Message: %s", topic, key, message)

	return nil
}
