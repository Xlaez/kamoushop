package routes

import (
	"kamoushop/pkg/controllers"
	"kamoushop/pkg/services/token"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(router *gin.Engine, c controllers.AuthController, token_maker token.Maker) {
	auth := router.Group("/v1/auth")
	auth.POST("/")
}
