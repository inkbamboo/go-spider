package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type BaseModel struct {
	ID        primitive.ObjectID `json:"id" bson:"_id" `
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at" `
}
