package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vadgun/gotrelloclone/notification-service/models"
	"github.com/vadgun/gotrelloclone/notification-service/services"
)

type NotificationHandler struct {
	service          *services.NotificationService
	webSocketService *services.WebSocketService
}

func NewNotificationHandler(service *services.NotificationService, wsService *services.WebSocketService) *NotificationHandler {
	return &NotificationHandler{service: service, webSocketService: wsService}
}

func (h *NotificationHandler) SendNotification(ctx *gin.Context) {
	var notification models.Notification
	if err := ctx.ShouldBindJSON(&notification); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	notification.CreatedAt = time.Now()
	err := h.service.CreateNotification(ctx, &notification)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo enviar la notificación"})
		return
	}

	// Enviar notificación a WebSocket service
	h.service.SendNotificationWebSocket(notification.Message)

	ctx.JSON(http.StatusOK, gin.H{"message": "Notificación enviada"})
}

// Método para manejar conexiones WebSocket service
func (h *NotificationHandler) HandleConnections(c *gin.Context) {
	h.webSocketService.HandleConnections(c.Writer, c.Request)
}
