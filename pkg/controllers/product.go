package controllers

import (
	"errors"
	"kamoushop/pkg/libs"
	"kamoushop/pkg/models"
	"kamoushop/pkg/services/api"
	"kamoushop/pkg/services/token"
	"kamoushop/pkg/services/types"
	"kamoushop/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProductController interface {
	CreateProduct() gin.HandlerFunc
	GetProductsByUserId() gin.HandlerFunc
	QueryProductsByName() gin.HandlerFunc
	GetProdById() gin.HandlerFunc
	DeleteProduct() gin.HandlerFunc
	UpdateProduct() gin.HandlerFunc
}

type productController struct {
	s      api.ProductService
	maker  token.Maker
	config utils.Config
}

func NewProductController(s api.ProductService, maker token.Maker, config utils.Config) ProductController {
	return &productController{
		s:      s,
		maker:  maker,
		config: config,
	}
}

func (p *productController) CreateProduct() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request types.Product
		if err := ctx.ShouldBind(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		secure_url, _, err := libs.UploadToCloud(ctx)
		if err != nil {
			ctx.JSON(http.StatusExpectationFailed, errorRes(err))
			return
		}

		payload := ctx.MustGet(authPayload).(*token.Payload)

		data := types.Product{
			Price:       request.Price,
			Name:        request.Name,
			Image:       secure_url,
			Description: request.Description,
			TotalStock:  request.TotalStock,
		}

		result, err := p.s.CreateProduct(data, payload.UserID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		ctx.JSON(http.StatusCreated, result)
	}
}

func (p *productController) GetProductsByUserId() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request types.GetProductsByUserId
		if err := ctx.ShouldBindQuery(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		user_id, err := primitive.ObjectIDFromHex(request.UserID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		counter := int64(1)
		skip := (request.Page - counter) * request.Limit
		filter := bson.D{{Key: "userId", Value: user_id}}
		options := &options.FindOptions{
			Limit: &request.Limit,
			Skip:  &skip,
		}
		result, totalDocs, err := p.s.GetProducts(filter, options)

		if err != nil {
			ctx.JSON(http.StatusExpectationFailed, errorRes(err))
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"products": result, "totalDocuments": totalDocs})
	}
}

func (p *productController) QueryProductsByName() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request types.QueryBrands
		if err := ctx.ShouldBindQuery(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		counter := int64(1)
		skip := (request.Page - counter) * request.Limit

		regex := primitive.Regex{Pattern: request.Keyword, Options: "i"}
		filter := bson.D{{Key: "name", Value: regex}}

		var products []models.Product

		products, totalDocs, err := p.s.GetProducts(filter, &options.FindOptions{Limit: &request.Limit, Skip: &skip})

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"products: ": products, "totalDocuments: ": totalDocs})
	}
}

func (p *productController) GetProdById() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request types.GetProdById
		if err := ctx.ShouldBindUri(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		id, err := primitive.ObjectIDFromHex(request.ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		var product models.Product
		product, err = p.s.GetProdById(id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		ctx.JSON(http.StatusOK, product)
	}
}

func (p *productController) DeleteProduct() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request types.GetProdById
		if err := ctx.ShouldBindUri(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		id, err := primitive.ObjectIDFromHex(request.ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		if err = p.s.DeleteProduct(id); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		ctx.JSON(http.StatusNoContent, msgRes(""))
	}
}

func (p *productController) UpdateProduct() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request types.UpdateProduct
		if err := ctx.ShouldBindJSON(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		id, err := primitive.ObjectIDFromHex(request.ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		filter := bson.D{primitive.E{Key: "_id", Value: id}}
		descrObj := bson.D{{Key: "description", Value: request.Description}}
		priceObj := bson.D{{Key: "price", Value: request.Price}}

		var updateObj bson.D

		if len(request.Description) > 1 && request.Price > 1 {
			updateObj = bson.D{{Key: "$set", Value: descrObj}, {Key: "$set", Value: priceObj}}
		} else if len(request.Description) > 1 {
			updateObj = bson.D{{Key: "$set", Value: descrObj}}
		} else if request.Price > 1 {
			updateObj = bson.D{{Key: "$set", Value: priceObj}}
		} else {
			ctx.JSON(http.StatusBadRequest, errorRes(errors.New("please provide a field to update")))
			return
		}

		if err = p.s.UpdateOne(filter, updateObj); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		ctx.JSON(http.StatusOK, msgRes("updated"))
	}
}
