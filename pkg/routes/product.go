package routes

import (
	"kamoushop/pkg/controllers"
	"kamoushop/pkg/middlewares"
	"kamoushop/pkg/services/token"

	"github.com/gin-gonic/gin"
)

func PoductRoutes(router *gin.Engine, c controllers.ProductController, token_maker token.Maker) {
	products := router.Group("/v1/product").Use(middlewares.AuthMiddleWare(token_maker))
	products.GET("/:id", c.GetProdById())
	products.GET("/products/by-id", c.GetProductsByUserId())
	products.GET("/products/by-name", c.QueryProductsByName())
	products.PATCH("/update", c.UpdateProduct())
	products.POST("/", c.CreateProduct())
	products.DELETE("/:id", c.DeleteProduct())
	products.POST("/add-to-cart", c.AddToCart())
}
