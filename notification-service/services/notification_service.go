package services

import (
	"context"
	"log"

	"github.com/vadgun/gotrelloclone/notification-service/models"
	"github.com/vadgun/gotrelloclone/notification-service/repositories"
)

type NotificationService struct {
	repo             *repositories.NotificationRepository
	webSocketService *WebSocketService
}

func NewNotificationService(repo *repositories.NotificationRepository, wsService *WebSocketService) *NotificationService {
	return &NotificationService{repo: repo, webSocketService: wsService}
}

func (s *NotificationService) CreateNotification(ctx context.Context, notification *models.Notification) error {
	return s.repo.SaveNotification(ctx, notification)
}

func (s *NotificationService) SendNotificationWebSocket(message string) {
	log.Println("Enviando notificación por WebSocket:", message)
	s.webSocketService.SendMessage(message)
}
