package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskStatus string

const (
	TODO       TaskStatus = "TODO"
	INPROGRESS TaskStatus = "IN_PROGRESS"
	DONE       TaskStatus = "DONE"
)

type Task struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title       string             `bson:"title" json:"title" binding:"required"`
	Description string             `bson:"description" json:"description" binding:"required"`
	BoardID     string             `bson:"board_id" json:"board_id" binding:"required"`
	UserID      string             `bson:"user_id" json:"user_id"`
	AssigneeID  string             `bson:"assignee_id,omitempty" json:"assignee_id,omitempty"`
	Status      TaskStatus         `bson:"status" json:"status"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}
