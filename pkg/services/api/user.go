package api

import (
	"context"
	"kamoushop/pkg/models"
	"kamoushop/pkg/services/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserService interface {
	GetUserById(id primitive.ObjectID) (types.User, error)
	GetUserByIdWithPassword(id primitive.ObjectID) (models.User, error)
	UpdateUser(filter bson.D, updateObj bson.D) error
	FindOne(filter bson.D) (models.User, error)
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

func (a *userService) GetUserByIdWithPassword(id primitive.ObjectID) (models.User, error) {
	user := models.User{}
	filter := bson.D{{Key: "_id", Value: id}}

	if err := a.col.FindOne(a.ctx, filter).Decode(&user); err == mongo.ErrNoDocuments && err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (u *userService) UpdateUser(filter bson.D, updateObj bson.D) error {
	if _, err := u.col.UpdateOne(u.ctx, filter, updateObj, options.Update()); err != nil {
		return err
	}
	return nil
}

func (u *userService) FindOne(filter bson.D) (models.User, error) {
	var user models.User
	if err := u.col.FindOne(u.ctx, filter).Decode(&user); err != nil {
		return models.User{}, err
	}
	return user, nil
}
