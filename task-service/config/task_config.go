package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database
var JWTSecret string

type JWTClaims struct {
	jwt.RegisteredClaims
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

// InitConfig carga las variables de entorno y configura la base de datos
func InitConfig() {
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
	fmt.Println("✅ Conectado a MongoDB desde Task-service:", dbName)
}
