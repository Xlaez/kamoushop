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
}
