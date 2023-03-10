package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	Price       int64              `json:"price" bson:"price"`
	Image       string             `json:"image,omitempty" bson:"image"`
	Name        string             `json:"name,omitempty" bson:"name"`
	UserID      primitive.ObjectID `json:"user_id,omitempty" bson:"userId"`
	Description string             `json:"description,omitempty" bson:"description"`
	CreatedAT   time.Time          `json:"created_at" bson:"createdAt"`
	UpdatedAT   time.Time          `json:"updated_at" bson:"updatedAt"`
}
