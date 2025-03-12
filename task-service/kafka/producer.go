package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type TaskProducer struct {
	writer *kafka.Writer
}

func NewTaskProducer(brokers []string) *TaskProducer {
	writer := &kafka.Writer{
		Addr:  kafka.TCP(brokers...),
		Async: true, // Opcional: para enviar mensajes de forma asíncrona
	}
	return &TaskProducer{writer: writer}
}

func (p *TaskProducer) PublishTaskEvent(topic string, message []byte) error {
	return p.writer.WriteMessages(context.Background(), kafka.Message{
		Topic: topic,
		Value: message,
	})
}
