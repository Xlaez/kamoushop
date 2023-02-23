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
	GetAllUsers(limit int64, page int64) ([]types.User, int64, error)
	QueryBrands(brand_name_keyword string, limit int64, page int64) ([]types.User, int64, error)
	DeleteUser(userId primitive.ObjectID) error
	// AddToCart(user_id primitive.ObjectID, cart []models.UserProduct) error
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

func (u *userService) GetAllUsers(limit int64, page int64) ([]types.User, int64, error) {
	counter := int64(1)
	skip := (page - counter) * limit
	cursor, err := u.col.Find(u.ctx, bson.D{}, &options.FindOptions{Limit: &limit, Skip: &skip})

	if err != nil {
		return nil, 0, err
	}

	count, err := u.col.CountDocuments(u.ctx, bson.D{})
	if err != nil {
		return nil, 0, err
	}

	var users []types.User

	if err = cursor.All(u.ctx, &users); err != nil {
		return nil, 0, err
	}

	return users, count, nil
}

func (u *userService) QueryBrands(brand_name_keyword string, limit int64, page int64) ([]types.User, int64, error) {
	counter := int64(1)
	skip := (page - counter) * limit
	filter := primitive.Regex{Pattern: brand_name_keyword, Options: "i"}
	cursor, err := u.col.Find(u.ctx, bson.D{{Key: "brandName", Value: filter}}, &options.FindOptions{Limit: &limit, Skip: &skip})

	if err != nil {
		return nil, 0, err
	}

	count, err := u.col.CountDocuments(u.ctx, bson.D{{Key: "brandName", Value: filter}})
	if err != nil {
		return nil, 0, err
	}

	var brands []types.User

	if err = cursor.All(u.ctx, &brands); err != nil {
		return nil, 0, err
	}

	return brands, count, nil
}

func (u *userService) DeleteUser(userId primitive.ObjectID) error {
	filter := bson.D{primitive.E{Key: "_id", Value: userId}}

	_, err := u.col.DeleteOne(u.ctx, filter, options.Delete())
	if err != nil {
		return err
	}

	return nil
}
