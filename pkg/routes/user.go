package routes

import (
	"kamoushop/pkg/controllers"
	"kamoushop/pkg/middlewares"
	"kamoushop/pkg/services/token"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine, c controllers.UserController, token_maker token.Maker) {
	user := router.Group("/v1/user").Use(middlewares.AuthMiddleWare(token_maker))
	user.GET("/", c.GetAllUsers())
	user.GET("/by-id/:id", c.GetUserById())
	user.GET("/me", c.GetUser())
	user.GET("/brands", c.QueryBrands())
	user.PATCH("/update/password", c.ChangePassword())
	user.PATCH("/update/image", c.UpdateImage())
	user.PATCH("/update/profile", c.UpdateProfile())
	user.PATCH("/update/brand-name", c.UpdateBrandName())
	user.PATCH("/star/:id", c.StarUserShop())
	user.DELETE("/:password", c.DeleteUser())
}
