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
	Stars      int64                `json:"stars,omitempty" bson:"stars" default:"0"`
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

type ValidateAcc struct {
	Code string `json:"code" binding:"required,min=6,max=6"`
}

type ChangePassword struct {
	OldPassword string `json:"old_password" binding:"required,min=7,alphanum"`
	NewPassword string `json:"new_password" binding:"required,min=7,alphanum"`
}

type UploadImage struct {
	ID string `form:"id" binding:"required"`
}

type UpdateProfile struct {
	PhoneNO   string `json:"phone_no"`
	Instagram string `json:"instagram"`
	Facebook  string `json:"facebook"`
}

type UpdateBrandName struct {
	BrandName string `json:"brand_name"`
}

type GetUsers struct {
	Limit int64 `form:"limit" biniding:"required"`
	Page  int64 `form:"page" binding:"required"`
}

type QueryBrands struct {
	Limit   int64  `form:"limit" biniding:"required"`
	Page    int64  `form:"page" binding:"required"`
	Keyword string `form:"keyword" binding:"required"`
}

type Product struct {
	Price       int    `form:"price" binding:"required"`
	Name        string `form:"name" binding:"required,min=3"`
	Image       string `form:"image"`
	Description string `form:"description" binding:"required,min=5"`
}

type GetProductsByUserId struct {
	Limit  int64  `form:"limit" biniding:"required"`
	Page   int64  `form:"page" binding:"required"`
	UserID string `form:"user_id" binding:"required"`
}

type DeleteUser struct {
	Password string `uri:"password" binding:"required"`
}

type StarShop struct {
	ID string `uri:"id" binding:"required"`
}

type GetProdById struct {
	ID string `uri:"id" binding:"required"`
}

type UpdateProduct struct {
	ID          string `json:"id" binding:"required"`
	Description string `json:"description"`
	Price       int64  `json:"price"`
}

type AddToCart struct {
	ProdID string `json:"prod_id" binidng:"required"`
}
