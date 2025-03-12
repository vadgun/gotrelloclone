package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/vadgun/gotrelloclone/task-service/config"
	"github.com/vadgun/gotrelloclone/task-service/handlers"
	"github.com/vadgun/gotrelloclone/task-service/kafka"
	"github.com/vadgun/gotrelloclone/task-service/repositories"
	"github.com/vadgun/gotrelloclone/task-service/routes"
	"github.com/vadgun/gotrelloclone/task-service/services"
)

func main() {
	// Iniciar conexión a MongoDB
	config.InitConfig()

	// Inicializar repositorio y servicio
	taskRepo := repositories.NewTaskRepository()
	kafkaProducer := kafka.NewTaskProducer([]string{"kafka:9092"})
	taskService := services.NewTaskService(taskRepo, kafkaProducer)
	taskHandler := handlers.NewTaskHandler(taskService)

	// Configurar el servicio en modo producción
	gin.SetMode(gin.ReleaseMode)

	// Configurar router Gin
	router := gin.Default()

	// Configurar rutas
	routes.SetupTaskRoutes(router, taskHandler)

	// Iniciar servidor en el puerto 8082
	log.Println("🚀 task-service corriendo en http://task-service:8080")
	router.Run(":8080")
}
