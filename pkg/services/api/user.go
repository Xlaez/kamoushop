package api

import (
	"context"
	"kamoushop/pkg/services/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserService interface {
	GetUser(id primitive.ObjectID) types.User
}

type userService struct {
	col *mongo.Collection
	ctx context.Context
}

func NewUserService(col *mongo.Collection, ctx context.Context) UserService {
	return &userService{
		col: col,
		ctx: ctx,
	}
}

func (u *userService) GetUser(id primitive.ObjectID) types.User {
	return types.User{}
}
