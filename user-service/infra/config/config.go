package config

import (
	"context"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"github.com/vadgun/gotrelloclone/user-service/infra/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
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
		logger.Log.Error("Faltan variables de entorno necesarias para mongo")
	}

	// Configurar conexión a MongoDB
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		logger.Log.Error("Error conectando a MongoDB:", zap.String("error", err.Error()))
	}

	// Verificar conexión
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = client.Ping(ctx, nil)
	if err != nil {
		logger.Log.Error("No se pudo conectar a MongoDB:", zap.String("error", err.Error()))
	}

	// Asignar base de datos
	DB = client.Database(dbName)
	logger.Log.Info("✅ Conectado a MongoDB desde user-service:", zap.String("db_Name", dbName))
}

func InitRedis() {
	// Obtener valores de entorno
	redisAddr := os.Getenv("REDIS_ADDR") // "localhost:6379" o nombre del contenedor Docker
	redisPass := os.Getenv("REDIS_PASS")

	// Verificar variables de entorno necesarias
	if redisAddr == "" {
		logger.Log.Error("Falta la variable de entorno REDIS_ADDR")
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
		logger.Log.Error("No se pudo conectar a Redis:", zap.String("error", err.Error()))
	}

	logger.Log.Info("✅ Conectado a Redis en:", zap.String("redis_Addr", redisAddr))
}
