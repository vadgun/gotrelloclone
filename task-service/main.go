package main

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/vadgun/gotrelloclone/task-service/config"
	"github.com/vadgun/gotrelloclone/task-service/handlers"
	"github.com/vadgun/gotrelloclone/task-service/kafka"
	"github.com/vadgun/gotrelloclone/task-service/logger"
	"github.com/vadgun/gotrelloclone/task-service/metrics"
	"github.com/vadgun/gotrelloclone/task-service/repositories"
	"github.com/vadgun/gotrelloclone/task-service/routes"
	"github.com/vadgun/gotrelloclone/task-service/services"
)

func main() {
	// Iniciar conexión a MongoDB
	config.InitConfig()

	// Inicializar metricas en Prometheus
	metrics.InitMetrics()

	// Iniciar el logger
	logger.InitLogger()

	// Inicializar repositorio y servicio
	taskRepo := repositories.NewTaskRepository()
	taskService := services.NewTaskService(taskRepo)
	taskHandler := handlers.NewTaskHandler(taskService)

	// Configurar el servicio en modo producción
	gin.SetMode(gin.ReleaseMode)

	// Configurar router Gin
	router := gin.Default()

	// Activando CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"POST", "OPTIONS", "GET", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Configurar rutas
	routes.SetupTaskRoutes(router, taskHandler)

	// Iniciar servidor en el puerto 8082
	logger.Log.Info("🚀 task-service corriendo en http://task-service:8080")
	go kafka.StartConsumer()
	router.Run(":8080")
	select {}
}
