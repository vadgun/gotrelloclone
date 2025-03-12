package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vadgun/gotrelloclone/notification-service/handlers"
)

func SetupNotificationRoutes(router *gin.Engine, notificationHandler *handlers.NotificationHandler) {
	router.POST("/notify", notificationHandler.SendNotification)
	router.GET("/ws", notificationHandler.HandleConnections)
}
