package kafka

import (
	"context"
	"encoding/json"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"
)

type KafkaProducer struct {
	Producer *kafka.Producer
	Topic    string
	Logger   *zap.Logger
}

func NewKafkaProducer(brokers string, topic string, logger *zap.Logger) *KafkaProducer {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": brokers,
	})
	if err != nil {
		logger.Fatal("❌ No se pudo crear el productor de Kafka", zap.Error(err))
	}

	go func() {
		for event := range producer.Events() {
			switch ev := event.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					logger.Error("❌ Error al enviar mensaje", zap.Error(ev.TopicPartition.Error))
				} else {
					logger.Info("📤 Mensaje entregado")
				}
			}
		}
	}()

	return &KafkaProducer{
		Producer: producer,
		Topic:    topic,
		Logger:   logger,
	}
}

func (kp *KafkaProducer) Publish(ctx context.Context, key string, value any) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		kp.Logger.Error("❌ Error al deserealizar el mensaje", zap.Error(err))
		return err
	}

	msg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &kp.Topic, Partition: kafka.PartitionAny},
		Key:            []byte(key),
		Value:          bytes,
	}

	err = kp.Producer.Produce(msg, nil)
	if err != nil {
		kp.Logger.Error("❌ Error al producir mensaje", zap.Error(err))
		return err
	}

	kp.Logger.Info("✅ Mensaje enviado a Kafka:", zap.String("key", key), zap.String("topic", kp.Topic))

	return nil
}

// // ProduceMessage envía un mensaje o publica el evento en Kafka
// func AProduceMessage(userID, message, topic, key string) error {
// 	p, err := kafka.NewProducer(&kafka.ConfigMap{
// 		"bootstrap.servers": "kafka:9092",
// 	})
// 	if err != nil {
// 		return fmt.Errorf("❌ Error al crear productor de Kafka: %v", err)
// 	}
// 	defer p.Close()

// 	// Enviar el mensaje o publicamos el evento en Kafka
// 	err = p.Produce(&kafka.Message{
// 		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
// 		Key:            []byte(key),
// 		Value:          []byte(message),
// 	}, nil)

// 	if err != nil {
// 		return fmt.Errorf("❌ Error enviando mensajea Kafka: %w", err)
// 	}

// 	log.Printf("✅ Mensaje enviado a Kafka | Topic: %s | Key: %s | Message: %s", topic, key, message)

// 	return nil
// }
