package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Board struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `json:"name" bson:"name"`
	OwnerID   string             `json:"owner_id" bson:"owner_id"`
	OwnerName string             `json:"owner_name" bson:"owner_name"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}
