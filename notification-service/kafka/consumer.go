package kafka

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// Task representa el esquema de la tarea recibida
type Task struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	UserID      string `json:"user_id"`
}

func sendNotification(task Task) {
	fmt.Printf("📢 Enviando notificación al usuario %s: Nueva tarea creada - %s\n", task.UserID, task.Title)
}

// StartConsumer inicia el consumidor de Kafka
func StartConsumer() {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "kafka:9092",
		"group.id":          "notification-service",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		log.Fatalf("❌ Error creando consumidor: %v", err)
	}
	defer c.Close()

	// Suscribirse al tópico
	topic := "task-events"
	err = c.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		log.Fatalf("❌ Error suscribiéndose a Kafka: %v", err)
	}

	log.Println("📩 Escuchando eventos de tareas en Kafka...")

	// Loop infinito para escuchar eventos
	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			log.Printf("📨 Notificación recibida | Topic: %s | Message: %s\n", *msg.TopicPartition.Topic, string(msg.Value))
			// Parsear JSON del mensaje
			var task Task
			if err := json.Unmarshal(msg.Value, &task); err != nil {
				log.Printf("⚠️ Error parseando JSON: %v\n", err)
				continue
			}

			// Aquí podrías llamar una función para enviar notificaciones
			sendNotification(task)
		} else {
			log.Printf("⚠️ Error al recibir mensaje: %v\n", err)
		}
	}
}
