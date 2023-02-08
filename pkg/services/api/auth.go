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
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	CreateUser(data models.User) error
	Login(data types.Login) (models.User, error)
	GetUserById(id string) (models.User, error)
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

func (a *authService) Login(data types.Login) (models.User, error) {
	user, err := GetUserByEmail(a, data.Email)

	if err != nil {
		return models.User{}, err
	}

	if err = password.ComparePassword(data.Password, user.Password); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			new_err := errors.New("password does not match")
			return models.User{}, new_err
		}
		return models.User{}, err
	}

	return user, nil
}

func GetUserByEmail(a *authService, email string) (models.User, error) {
	user := models.User{}
	filter := bson.D{{Key: "email", Value: email}}

	if err := a.col.FindOne(a.ctx, filter).Decode(&user); err == mongo.ErrNoDocuments && err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (a *authService) GetUserById(id string) (models.User, error) {
	user_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.User{}, err
	}

	user := models.User{}
	filter := bson.D{{Key: "_id", Value: user_id}}

	if err := a.col.FindOne(a.ctx, filter).Decode(&user); err == mongo.ErrNoDocuments && err != nil {
		return models.User{}, err
	}
	return user, nil
}
