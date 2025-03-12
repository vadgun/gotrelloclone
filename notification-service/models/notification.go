package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Notification struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    string             `bson:"user_id" json:"user_id" binding:"required"`
	Message   string             `bson:"message" json:"message" binding:"required"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
}
