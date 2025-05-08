package main

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/vadgun/gotrelloclone/user-service/handlers"
	"github.com/vadgun/gotrelloclone/user-service/infra/config"
	"github.com/vadgun/gotrelloclone/user-service/infra/kafka"
	"github.com/vadgun/gotrelloclone/user-service/infra/logger"
	"github.com/vadgun/gotrelloclone/user-service/infra/metrics"
	"github.com/vadgun/gotrelloclone/user-service/middlewares"
	"github.com/vadgun/gotrelloclone/user-service/repositories"
	"github.com/vadgun/gotrelloclone/user-service/routes"
	"github.com/vadgun/gotrelloclone/user-service/services"
)

func main() {
	// Inicializar el logger
	logger.InitLogger()
	log := logger.Log

	// Iniciar conexión a MongoDB
	config.InitMongo()

	// Iniciar conexión a Redis
	config.InitRedis()

	// Inicializar metricas en Prometheus
	metrics.InitMetrics()

	// Inicializar repositorio
	userRepo := repositories.NewUserRepository(
		config.DB.Collection("users"),
		config.RedisClient,
		log,
	)

	// Productor de kafka
	kafkaProducer := kafka.NewKafkaProducer("kafka:9092", "user-events", log)

	// Inicializar el servicio
	userService := services.NewUserService(userRepo, log, kafkaProducer)

	// Inicializar el handler
	userHandler := handlers.NewUserHandler(userService, log)

	// Configurar el servicio en modo producción
	gin.SetMode(gin.ReleaseMode)

	// Configurar router Gin
	router := gin.Default()

	// Acivando CORS para nuestro FrontEnd
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"POST", "OPTIONS", "GET", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Configurar rutas
	routes.SetupUserRoutes(router, userHandler)

	// Limitar la tasa de peticiones
	router.Use(middlewares.RateLimitMiddleware())

	// Envolver el manejador de Prometheus/http para rutearlo a gin
	router.Use(metrics.MetricsMiddleware())
	router.GET("/metrics", gin.WrapH(metrics.MetricsHandler()))

	// Iniciar servidor en el puerto 8080
	log.Info("🚀 user-service corriendo en http://user-service:8080")
	router.Run(":8080")
}
