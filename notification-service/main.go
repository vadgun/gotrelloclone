package main

import (
	"github.com/gin-gonic/gin"
	"github.com/vadgun/gotrelloclone/notification-service/config"
	"github.com/vadgun/gotrelloclone/notification-service/handlers"
	"github.com/vadgun/gotrelloclone/notification-service/kafka"
	"github.com/vadgun/gotrelloclone/notification-service/logger"
	"github.com/vadgun/gotrelloclone/notification-service/metrics"
	"github.com/vadgun/gotrelloclone/notification-service/repositories"
	"github.com/vadgun/gotrelloclone/notification-service/routes"
	"github.com/vadgun/gotrelloclone/notification-service/services"
)

func main() {
	config.InitConfig()
	logger.InitLogger()
	metrics.InitMetrics()

	notificationRepo := repositories.NewNotificationRepository()
	webSocketService := services.NewWebSocketService()
	go webSocketService.HandleMessages()
	notificationService := services.NewNotificationService(notificationRepo, webSocketService)
	notificationHandler := handlers.NewNotificationHandler(notificationService, webSocketService)

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	routes.SetupNotificationRoutes(router, notificationHandler)
	router.GET("/metrics", gin.WrapH(metrics.MetricsHandler()))
	logger.Log.Info("🚀 notification-service corriendo en http://notification-service:8080")
	go kafka.StartConsumer(notificationHandler)
	router.Run(":8080")
	select {}
}
