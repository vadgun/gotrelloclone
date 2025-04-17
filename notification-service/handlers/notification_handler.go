package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vadgun/gotrelloclone/notification-service/logger"
	"github.com/vadgun/gotrelloclone/notification-service/metrics"
	"github.com/vadgun/gotrelloclone/notification-service/models"
	"github.com/vadgun/gotrelloclone/notification-service/services"
	"go.uber.org/zap"
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

	// Incrementar la métrica cada vez que se llame este endpoint
	metrics.HttpRequestsTotal.WithLabelValues("POST", "/notify").Inc()

	notification.CreatedAt = time.Now()
	err := h.service.CreateNotification(ctx, &notification)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo enviar la notificación"})
		return
	}

	// Enviar notificación a WebSocket service
	h.service.SendNotificationWebSocket(notification.Message)

	// Creando un log personalizado cuando se crea una notificacion
	logger.Log.Info("Creando Notificacion", zap.String("endpoint", ctx.Request.URL.Path), zap.String("ip", ctx.ClientIP()))

	ctx.JSON(http.StatusOK, gin.H{"message": "Notificación enviada"})
}

// Método para manejar conexiones WebSocket service
func (h *NotificationHandler) HandleConnections(c *gin.Context) {
	h.webSocketService.HandleConnections(c.Writer, c.Request)
}

func (h *NotificationHandler) HandleKafkaMessage(notification models.Notification) error {
	// var notification models.Notification
	// if err := json.Unmarshal(message.Value, &notification); err != nil {
	// 	return fmt.Errorf("error al deserializar el mensaje: %v", err)
	// }

	// notification.CreatedAt = time.Now()

	// Guardar la notificación directamente usando el repositorio
	if err := h.service.CreateNotification(context.Background(), &notification); err != nil {
		return fmt.Errorf("error al guardar la notificación: %v", err)
	}

	// Enviar notificación a WebSocket service
	h.service.SendNotificationWebSocket(notification.Message)

	return nil
}
