package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	DB          *mongo.Database
	RedisClient *redis.Client
	JWTSecret   string
)

type JWTClaims struct {
	jwt.RegisteredClaims
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

// InitConfig carga las variables de entorno y configura la base de datos
func InitMongo() {
	// Obtener valores de entorno
	mongoURI := os.Getenv("MONGO_URI")
	dbName := os.Getenv("MONGO_DB_NAME")
	JWTSecret = os.Getenv("JWT_SECRET")

	if mongoURI == "" || dbName == "" || JWTSecret == "" {
		log.Fatal("Faltan variables de entorno necesarias")
	}

	// Configurar conexión a MongoDB
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal("Error conectando a MongoDB:", err)
	}

	// Verificar conexión
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("No se pudo conectar a MongoDB:", err)
	}

	// Asignar base de datos
	DB = client.Database(dbName)
	fmt.Println("✅ Conectado a MongoDB desde user-service:", dbName)
}

func InitRedis() {
	// Obtener valores de entorno
	redisAddr := os.Getenv("REDIS_ADDR") // "localhost:6379" o nombre del contenedor Docker
	redisPass := os.Getenv("REDIS_PASS")

	// Verificar variables de entorno necesarias
	if redisAddr == "" {
		log.Fatal("Falta la variable de entorno REDIS_ADDR")
	}

	// Inicializar cliente Redis
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPass,
		DB:       0,
	})

	// Verificar conexión
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatal("No se pudo conectar a Redis:", err)
	}

	fmt.Println("✅ Conectado a Redis en:", redisAddr)
}
