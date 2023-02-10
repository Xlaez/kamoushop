package routes

import (
	"kamoushop/pkg/controllers"
	"kamoushop/pkg/middlewares"
	"kamoushop/pkg/services/token"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine, c controllers.UserController, token_maker token.Maker) {
	user := router.Group("/v1/user").Use(middlewares.AuthMiddleWare(token_maker))
	user.GET("/by-id/:id", c.GetUserById())
	user.GET("/me", c.GetUser())
	user.PATCH("/update/password", c.ChangePassword())
	user.PATCH("/update/image", c.UpdateImage())
	user.PATCH("/update/profile", c.UpdateProfile())
	user.PATCH("/update/brand-name", c.UpdateBrandName())
}
