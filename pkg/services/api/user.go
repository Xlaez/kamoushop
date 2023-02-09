package api

import (
	"context"
	"kamoushop/pkg/services/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserService interface {
	GetUserById(id primitive.ObjectID) (types.User, error)
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

func (a *userService) GetUserById(id primitive.ObjectID) (types.User, error) {
	user := types.User{}
	filter := bson.D{{Key: "_id", Value: id}}

	if err := a.col.FindOne(a.ctx, filter).Decode(&user); err == mongo.ErrNoDocuments && err != nil {
		return types.User{}, err
	}
	return user, nil
}
