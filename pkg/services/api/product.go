package api

import (
	"context"
	"fmt"
	"kamoushop/pkg/models"
	"kamoushop/pkg/services/types"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProductService interface {
	CreateProduct(prod types.Product, userId primitive.ObjectID) (*mongo.InsertOneResult, error)
	GetProducts(filter bson.D, options *options.FindOptions) ([]models.Product, int64, error)
	GetProdById(id primitive.ObjectID) (models.Product, error)
	DeleteProduct(id primitive.ObjectID) error
	UpdateOne(filter bson.D, updateObj bson.D) error
}

type productService struct {
	col *mongo.Collection
	ctx context.Context
}

func NewProductService(ctx context.Context, col *mongo.Collection) ProductService {
	return &productService{
		col: col,
		ctx: ctx,
	}
}

func (p *productService) CreateProduct(prod types.Product, userId primitive.ObjectID) (*mongo.InsertOneResult, error) {
	id := primitive.NewObjectID()

	product := models.Product{
		ID:          id,
		Price:       fmt.Sprint("$", prod.Price),
		Image:       prod.Image,
		Name:        prod.Name,
		Description: prod.Description,
		TotalStock:  prod.TotalStock,
		UserID:      userId,
		CreatedAT:   time.Now(),
		UpdatedAT:   time.Now(),
	}

	result, err := p.col.InsertOne(p.ctx, &product, options.InsertOne())
	if err != nil {
		return &mongo.InsertOneResult{}, err
	}

	return result, nil
}

func (p *productService) GetProducts(filter bson.D, options *options.FindOptions) ([]models.Product, int64, error) {
	products := []models.Product{}
	cursor, err := p.col.Find(p.ctx, filter, options)

	if err != nil {
		return nil, 0, err
	}

	if err = cursor.All(p.ctx, &products); err != nil {
		return nil, 0, err
	}

	count, err := p.col.CountDocuments(p.ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return products, count, nil
}

func (p *productService) GetProdById(id primitive.ObjectID) (models.Product, error) {
	var product models.Product
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	if err := p.col.FindOne(p.ctx, filter, options.FindOne()).Decode(&product); err != nil {
		return models.Product{}, err
	}
	return product, nil
}

func (p *productService) DeleteProduct(id primitive.ObjectID) error {
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	if _, err := p.col.DeleteOne(p.ctx, filter, options.Delete()); err != nil {
		return err
	}
	return nil
}

func (p *productService) UpdateOne(filter bson.D, updateObj bson.D) error {
	if _, err := p.col.UpdateOne(p.ctx, filter, updateObj, options.Update()); err != nil {
		return err
	}
	return nil
}
