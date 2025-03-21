package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/vadgun/gotrelloclone/user-service/config"
	"github.com/vadgun/gotrelloclone/user-service/handlers"
	"github.com/vadgun/gotrelloclone/user-service/logger"
	"github.com/vadgun/gotrelloclone/user-service/metrics"
	"github.com/vadgun/gotrelloclone/user-service/repositories"
	"github.com/vadgun/gotrelloclone/user-service/routes"
	"github.com/vadgun/gotrelloclone/user-service/services"
)

func main() {
	// Iniciar conexión a MongoDB
	config.InitConfig()

	// Inicializar metricas en Prometheus
	metrics.InitMetrics()

	// Inicializar el logger
	logger.InitLogger()

	// Inicializar repositorio y servicio
	userRepo := repositories.NewUserRepository()
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	// Configurar el servicio en modo producción
	gin.SetMode(gin.ReleaseMode)

	// Configurar router Gin
	router := gin.Default()

	// Configurar rutas
	routes.SetupUserRoutes(router, userHandler)

	// Envolver el manejador de Prometheus/http para rutearlo a gin
	router.GET("/metrics", gin.WrapH(metrics.MetricsHandler()))

	// Iniciar servidor en el puerto 8080
	logrus.Info("🚀 user-service corriendo en http://user-service:8080")
	router.Run(":8080")
}
