package main

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/vadgun/gotrelloclone/user-service/handlers"
	"github.com/vadgun/gotrelloclone/user-service/infra/config"
	"github.com/vadgun/gotrelloclone/user-service/infra/logger"
	"github.com/vadgun/gotrelloclone/user-service/infra/metrics"
	"github.com/vadgun/gotrelloclone/user-service/repositories"
	"github.com/vadgun/gotrelloclone/user-service/routes"
	"github.com/vadgun/gotrelloclone/user-service/services"
)

func main() {
	// Inicializar el logger
	logger.InitLogger()

	// Iniciar conexión a MongoDB
	config.InitMongo()

	// Iniciar conexión a Redis
	config.InitRedis()

	// Inicializar metricas en Prometheus
	metrics.InitMetrics()

	// Inicializar repositorio y servicio
	userRepo := repositories.NewUserRepository()
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

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

	// Envolver el manejador de Prometheus/http para rutearlo a gin
	router.Use(metrics.MetricsMiddleware())
	router.GET("/metrics", gin.WrapH(metrics.MetricsHandler()))

	// Iniciar servidor en el puerto 8080
	logger.Log.Info("🚀 user-service corriendo en http://user-service:8080")
	router.Run(":8080")
}
