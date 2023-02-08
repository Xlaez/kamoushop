package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID         primitive.ObjectID   `json:"id,omitempty" bson:"_id"`
	FirstName  string               `json:"first_name,omitempty" bson:"firstname"`
	LastName   string               `json:"last_name,omitempty" bson:"lastname"`
	Image      string               `json:"image,omitempty" bson:"image"`
	BrandName  string               `json:"brand_name,omitempty" bson:"brandName"`
	PhoneNO    string               `json:"phone_no,omitempty" bson:"phoneNo"`
	Email      string               `json:"email,omitempty" bson:"email"`
	Instagram  string               `json:"instagram,omitempty" bson:"instagram"`
	Facebook   string               `json:"facebook,omitempty" bson:"facebook"`
	Stars      string               `json:"stars,omitempty" bson:"stars" default:"0"`
	StarredBy  []primitive.ObjectID `json:"starred_by" bson:"starredBy"`
	IsVerified bool                 `json:"is_verified" bson:"isVerified" default:"false"`
	CreatedAT  time.Time            `json:"created_at" bson:"createdAt"`
	UpdatedAT  time.Time            `json:"updated_at" bson:"updatedAt"`
}

type AddUser struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,alphanum,min=7"`
}

type Login struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,alphanum,min=7"`
}

type GetUser struct {
	ID string `uri:"id" binding:"required"`
}
