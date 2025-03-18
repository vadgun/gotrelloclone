package kafka

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/vadgun/gotrelloclone/notification-service/handlers"
	"github.com/vadgun/gotrelloclone/notification-service/models"
)

// Task representa el esquema de la tarea recibida mediante el evento de Kafka
type Task struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	UserID      string `json:"user_id"`
}

// StartConsumer inicia el consumidor de Kafka en notification-service, escuchando eventos de task.
func StartConsumer(kafkaHandler *handlers.NotificationHandler) {
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

	log.Println("📩 Escuchando eventos de tareas (task-service) en Kafka...")

	// Loop infinito para escuchar eventos
	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {

			// Diferenciar eventos por tópico y por key para decidir que hacer en nuestra base de datos de notification-service
			switch *msg.TopicPartition.Topic {
			case "task-events":
				log.Printf("📨 Evento recibido | Topic: %s | Message: %s| Key: %s\n", *msg.TopicPartition.Topic, string(msg.Value), string([]byte(msg.Key)))
				switch string([]byte(msg.Key)) {
				case "new-task":
					var task Task
					if err := json.Unmarshal(msg.Value, &task); err != nil {
						log.Printf("⚠️ Error parseando JSON de tarea: %v\n", err)
						continue
					}

					// Crear la notificación a partir de la tarea
					notification := models.Notification{
						UserID:    task.UserID,
						Message:   fmt.Sprintf("Nueva tarea creada: %s", task.Title),
						CreatedAt: time.Now(),
					}

					// Guardar la notificación y enviarla a WebSocket
					if err := kafkaHandler.HandleKafkaMessage(notification); err != nil {
						log.Printf("⚠️ Error manejando la notificación desde Kafka: %v\n", err)
					}

				case "update-task-status":
					var task Task
					if err := json.Unmarshal(msg.Value, &task); err != nil {
						log.Printf("⚠️ Error parseando JSON de tarea: %v\n", err)
						continue
					}

					// Crear la notificación a partir de la tarea
					notification := models.Notification{
						UserID:    task.UserID,
						Message:   fmt.Sprintf("Estado de la tarea actualizado: %s", task.Title),
						CreatedAt: time.Now(),
					}

					// Guardar la notificación y enviarla a WebSocket
					if err := kafkaHandler.HandleKafkaMessage(notification); err != nil {
						log.Printf("⚠️ Error manejando la notificación: %v\n", err)
					}
				}
			}

		} else {
			log.Printf("⚠️ Error al recibir mensaje: %v\n", err)
		}
	}
}
