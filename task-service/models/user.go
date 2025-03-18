package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserFromUserService struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`
}
