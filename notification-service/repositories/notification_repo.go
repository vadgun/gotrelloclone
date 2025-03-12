package repositories

import (
	"context"

	"github.com/vadgun/gotrelloclone/notification-service/config"
	"github.com/vadgun/gotrelloclone/notification-service/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type NotificationRepository struct {
	collection *mongo.Collection
}

func NewNotificationRepository() *NotificationRepository {
	return &NotificationRepository{
		collection: config.DB.Collection("notifications"),
	}
}

func (r *NotificationRepository) SaveNotification(ctx context.Context, notification *models.Notification) error {
	_, err := r.collection.InsertOne(ctx, notification)
	return err
}
