package routes

import (
	"kamoushop/pkg/controllers"
	"kamoushop/pkg/services/token"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine, c controllers.UserController, token_maker token.Maker) {
	user := router.Group("/v1/user")
	user.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"msg": "reached user route"})
	})
}
