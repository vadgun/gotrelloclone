package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func CleanUpTopic(brokers, topic string) error {

	adminClient, err := kafka.NewAdminClient(&kafka.ConfigMap{"bootstrap.servers": brokers})
	if err != nil {
		return fmt.Errorf("❌ no se pudo crear admin Kafka: %w", err)
	}

	defer adminClient.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	results, err := adminClient.DeleteTopics(
		ctx,
		[]string{topic},
		kafka.SetAdminOperationTimeout(5*time.Second),
	)
	if err != nil {
		return fmt.Errorf("❌ error al borrar topic: %w", err)
	}

	for _, result := range results {
		if result.Error.Code() != kafka.ErrNoError {
			return fmt.Errorf("❌ fallo al borrar topic %s: %s", result.Topic, result.Error.String())
		}
	}

	return nil
}

func CreateTopic(brokers, topic string) error {

	adminClient, err := kafka.NewAdminClient(&kafka.ConfigMap{
		"bootstrap.servers": brokers,
	})
	if err != nil {
		return fmt.Errorf("❌ no se pudo crear admin Kafka: %w", err)
	}
	defer adminClient.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mytopic := &kafka.TopicSpecification{
		Topic:         topic,
		NumPartitions: 1,
	}

	results, err := adminClient.CreateTopics(
		ctx,
		[]kafka.TopicSpecification{
			*mytopic,
		},
		kafka.SetAdminOperationTimeout(5*time.Second),
	)
	if err != nil {
		return fmt.Errorf("❌ error al crear el topic: %w", err)
	}

	for _, result := range results {
		if result.Error.Code() == kafka.ErrTopicAlreadyExists {
			fmt.Printf("⚠️ El topic %s ya existe\n", result.Topic)
			continue
		}
		if result.Error.Code() != kafka.ErrNoError {
			return fmt.Errorf("❌ fallo al crear topic %s: %s", result.Topic, result.Error.String())
		}
	}

	return nil
}
