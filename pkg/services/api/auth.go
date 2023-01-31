package api

import (
	"context"
	"errors"
	"kamoushop/pkg/models"
	"kamoushop/pkg/services/password"
	"kamoushop/pkg/services/types"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthService interface {
	CreateUser(data models.User) error
}

type authService struct {
	col *mongo.Collection
	ctx context.Context
}

func NewAuthService(col *mongo.Collection, ctx context.Context) AuthService {
	return &authService{
		col: col,
		ctx: ctx,
	}
}

func (a *authService) CreateUser(data models.User) error {
	id := primitive.NewObjectID()
	hashedPass, err := password.HashPassword(data.Password)

	if err != nil {
		return err
	}

	if user, _ := GetUserByEmail(a, data.Email); user.Email == data.Email {
		err = errors.New("user already exists")
		return err
	}

	new_user := models.User{
		ID:        id,
		FirstName: data.FirstName,
		LastName:  data.LastName,
		Password:  hashedPass,
		Email:     data.Email,
		CreatedAT: time.Now(),
		UpdatedAT: time.Now(),
	}

	_, err = a.col.InsertOne(a.ctx, new_user)
	if err != nil {
		return err
	}
	return nil
}

func GetUserByEmail(a *authService, email string) (types.User, error) {
	var user types.User
	filter := bson.D{{Key: "email", Value: email}}

	if err := a.col.FindOne(a.ctx, filter).Decode(&user); err != mongo.ErrNoDocuments && err != nil {
		return types.User{}, err
	}
	return user, nil
}
