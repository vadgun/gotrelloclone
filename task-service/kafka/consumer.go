package kafka

import (
	"encoding/json"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/vadgun/gotrelloclone/task-service/models"
	"github.com/vadgun/gotrelloclone/task-service/repositories"
)

// StartConsumer inicia el consumidor de Kafka en task-service, escuchando eventos de user.
func StartConsumer() {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "kafka:9092",
		"group.id":          "task-service",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		log.Fatalf("❌ Error creando consumidor: %v", err)
	}
	defer c.Close()

	// Suscribirse al topic
	topic := "user-events"
	err = c.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		log.Fatalf("❌ Error suscribiéndose a Kafka: %v", err)
	}

	log.Println("📩 Escuchando eventos de user-service para task-service en Kafka...")

	// Loop infinito para escuchar eventos
	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			switch *msg.TopicPartition.Topic {
			case "user-events":
				log.Printf("📨 Evento recibido | Topic: %s | Message: %s| Key: %s\n", *msg.TopicPartition.Topic, string(msg.Value), string([]byte(msg.Key)))
				switch string([]byte(msg.Key)) {
				case "new-user":
					// Parsear JSON del mensaje
					var user models.User
					if err := json.Unmarshal(msg.Value, &user); err != nil {
						log.Printf("⚠️ Error parseando JSON de usuario: %v\n", err)
						continue
					}

					// Guardar usuario en task-mongo
					userRepo := repositories.NewTaskRepository()
					err = userRepo.SaveUser(&user)
					if err != nil {
						log.Printf("⚠️ Error guardando usuario en task-service: %v\n", err)
					} else {
						log.Printf("✅ Usuario almacenado en task-service: %v\n", user)
					}
				}
			}
		} else {
			log.Printf("⚠️ Error al recibir mensaje: %v\n", err)
		}
	}
}
