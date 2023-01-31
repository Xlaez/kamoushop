package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Token struct {
	Token  string             `json:"token" bson:"token" required:"true"`
	UserID primitive.ObjectID `json:"user_id" bson:"userId" required:"true"`
	// possible values include: ["refresh", "access"]
	Type        string        `json:"type" bson:"type" required:"true"`
	ExpiresAT   time.Duration `json:"expires_at" bson:"expiresAt"`
	BlackListed bool          `json:"black_listed" bson:"blackListed" default:"false"`
}
