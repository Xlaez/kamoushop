package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	ID         primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	UserID     primitive.ObjectID `json:"user_id" bson:"userId"`
	TotalPrice int64              `json:"total_price" bson:"totalPrice"`
	CreatedAT  time.Time          `json:"created_at" bson:"createdAt"`
	UpdatedAT  time.Time          `json:"updated_at" bson:"updatedAt"`
}
