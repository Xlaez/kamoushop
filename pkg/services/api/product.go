package api

import (
	"context"
	"errors"
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
	AddToCart(product_id primitive.ObjectID, user_id primitive.ObjectID) error
	RemoveFromCart(product_id primitive.ObjectID, user_id primitive.ObjectID) error
	MakeOrder(user_id primitive.ObjectID) (models.Order, error)
}

type productService struct {
	col       *mongo.Collection
	ctx       context.Context
	user_col  *mongo.Collection
	order_col *mongo.Collection
}

func NewProductService(ctx context.Context, col *mongo.Collection, user_col *mongo.Collection, order_col *mongo.Collection) ProductService {
	return &productService{
		col:       col,
		ctx:       ctx,
		user_col:  user_col,
		order_col: order_col,
	}
}

var (
	ErrCantFindProduct = errors.New("can't find product")
	ErrCantUpdateUser  = errors.New("cannot add product to cart")
	ErrCantRemoveItem  = errors.New("cannot remove item from cart")
	ErrCantGetItem     = errors.New("cannot get item from cart ")
)

func (p *productService) CreateProduct(prod types.Product, userId primitive.ObjectID) (*mongo.InsertOneResult, error) {
	id := primitive.NewObjectID()

	product := models.Product{
		ID:          id,
		Price:       int64(prod.Price),
		Image:       prod.Image,
		Name:        prod.Name,
		Description: prod.Description,
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

func (p *productService) AddToCart(product_id primitive.ObjectID, user_id primitive.ObjectID) error {
	cursor, err := p.col.Find(p.ctx, bson.D{primitive.E{Key: "_id", Value: product_id}})

	if err != nil {
		return ErrCantFindProduct
	}

	var cart_prod []models.Prod

	if err = cursor.All(p.ctx, &cart_prod); err != nil {
		return err
	}

	updateObj := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "userCart.products", Value: bson.D{{Key: "$each", Value: cart_prod}}}}}}

	filter := bson.D{primitive.E{Key: "_id", Value: user_id}}
	if _, err := p.user_col.UpdateOne(p.ctx, filter, updateObj); err != nil {
		return ErrCantUpdateUser
	}
	return nil
}

func (p *productService) RemoveFromCart(product_id primitive.ObjectID, user_id primitive.ObjectID) error {
	filter := bson.D{primitive.E{Key: "_id", Value: user_id}}
	update := bson.M{"$pull": bson.M{"userCart.products": bson.M{"_id": product_id}}}
	if _, err := p.user_col.UpdateMany(p.ctx, filter, update); err != nil {
		return err
	}
	return nil
}

func (p *productService) MakeOrder(user_id primitive.ObjectID) (models.Order, error) {
	var user models.User
	filter := bson.D{primitive.E{Key: "_id", Value: user_id}}

	if err := p.user_col.FindOne(p.ctx, filter, options.FindOne()).Decode(&user); err != nil {
		return models.Order{}, err
	}

	var price int64 = 0
	for _, prod := range user.UserCart.Products {
		price += prod.Price
	}

	order := models.Order{
		ID:         primitive.NewObjectID(),
		UserID:     user_id,
		TotalPrice: price,
		CreatedAT:  time.Now(),
		UpdatedAT:  time.Now(),
	}

	// Add to the order model
	_, err := p.order_col.InsertOne(p.ctx, order, options.InsertOne())
	if err != nil {
		return models.Order{}, err
	}

	user_empty_cart := models.UserCart{
		Products: []models.Prod{},
	}
	// Delete user cart
	updateObj := bson.D{{Key: "$set", Value: bson.D{{Key: "userCart", Value: user_empty_cart}}}}
	if _, err = p.user_col.UpdateByID(p.ctx, user_id, updateObj, options.Update()); err != nil {
		return models.Order{}, err
	}

	return order, nil
}
