package integration

import (
	"context"
	"encoding/json"
	"log"
	"testing"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	mk "github.com/vadgun/gotrelloclone/user-service/infra/kafka"
	"github.com/vadgun/gotrelloclone/user-service/models"
	"github.com/vadgun/gotrelloclone/user-service/repositories"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func TestIntegration_CreateUser(t *testing.T) {
	// Configurar conexión real a MongoDB
	// Conectarse a una base de datos real (puede ser una de pruebas como gotrelloclone_test)
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27020")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	assert.NoError(t, err)

	// Inicializar repositorio real
	userCollection := client.Database("gotrelloclone_test").Collection("users")
	repo := &repositories.UserRepository{Collection: userCollection}

	// Crear usuario real
	user := &models.User{
		Name:     "Usuario Integracion",
		Email:    "integra@example.com",
		Password: "hashedpassword123",
		Phone:    "123456789",
		Role:     "user",
	}

	// Ejecutar la lógica real
	id, err := repo.CreateUser(user)

	// Verificar resultados reales
	assert.NoError(t, err)
	assert.NotEmpty(t, id)
}

func TestIntegration_GetUserByEmail(t *testing.T) {

	// Logger
	logger, _ := zap.NewDevelopment()

	// Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Cambia si estás usando Docker o localhost
	})
	defer redisClient.Close()
	_, err := redisClient.Ping(context.TODO()).Result()
	assert.NoError(t, err)

	// Crear una conexion real a Mongo
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27020")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	assert.NoError(t, err)

	// Conectar a una coleccion real de Mongo usando el repositorio de user-service
	userCollection := client.Database("gotrelloclone_test").Collection("users")
	repo := &repositories.UserRepository{Collection: userCollection, RedisClient: redisClient, Logger: logger}

	email := "cachetest@example.com"
	user := &models.User{
		Name:     "Cache Test",
		Email:    email,
		Password: "secretagent",
		Phone:    "1111111111",
		Role:     "user",
	}

	// Limpiar antes por si ya existe
	_, _ = userCollection.DeleteMany(context.TODO(), map[string]any{"email": email})
	_ = redisClient.Del(context.TODO(), "user_email:"+email).Err()

	// Insertar usuario en Mongo
	id, err := repo.CreateUser(user)
	assert.NoError(t, err)
	assert.NotEmpty(t, id)

	// Obtener usuario (primera vez, debe venir de Mongo)
	foundUser, err := repo.GetUserByEmail(email)
	assert.NoError(t, err)
	assert.Equal(t, user.Email, foundUser.Email)

	// Verificar que ahora está en Redis
	cached, err := redisClient.Get(context.TODO(), "user_email:"+email).Result()
	assert.NoError(t, err)

	var cachedUser models.User
	err = json.Unmarshal([]byte(cached), &cachedUser)
	assert.NoError(t, err)
	assert.Equal(t, user.Email, cachedUser.Email)

	// Segunda vez (debería venir desde Redis si todo va bien)
	start := time.Now()
	cachedResult, err := repo.GetUserByEmail(email)
	elapsed := time.Since(start)

	assert.NoError(t, err)
	assert.Equal(t, user.Email, cachedResult.Email)

	t.Logf("Tiempo de respuesta desde cache (aprox): %s", elapsed)

}

func TestKafkaProducerAndConsumer(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	topic := "user_created_test"
	brokers := "localhost:9092"

	// 1. Crear el productor
	producer := mk.NewKafkaProducer(brokers, topic, logger)

	// 2. Crear un mensaje de prueba
	type UserCreatedEvent struct {
		ID string `json:"id"`
	}

	event := UserCreatedEvent{
		ID: "test123",
	}

	// 4. Crear un consumidor temporal
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": brokers,
		"group.id":          "test-group",
		"auto.offset.reset": "earliest",
	})
	assert.NoError(t, err)
	defer consumer.Close()

	err = consumer.SubscribeTopics([]string{topic}, nil)
	assert.NoError(t, err)

	// 3. Publicar el evento
	ctx := context.Background()
	err = producer.Publish(ctx, event.ID, event)
	assert.NoError(t, err)

	// 5. Esperar y leer el mensaje
	msg, err := consumer.ReadMessage(25 * time.Second)
	if err == nil {
		switch *msg.TopicPartition.Topic {
		case "user_created_test":
			log.Printf("📨 Evento recibido | Topic: %s | Message: %s| Key: %s\n", *msg.TopicPartition.Topic, string(msg.Value), string([]byte(msg.Key)))
		}
	} else {
		log.Printf("⚠️ Error al recibir mensaje: %v\n", err)
	}
	assert.NoError(t, err)
	assert.NotEmpty(t, msg)

	// 6. Validar contenido
	var receivedEvent UserCreatedEvent
	err = json.Unmarshal(msg.Value, &receivedEvent)
	assert.NoError(t, err)

	assert.Equal(t, event.ID, receivedEvent.ID)
}
