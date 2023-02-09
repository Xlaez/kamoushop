package controllers

import (
	"kamoushop/pkg/services/api"
	"kamoushop/pkg/services/token"
	"kamoushop/pkg/services/types"
	"kamoushop/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	authPayload = "x-auth-token_payload"
)

type UserController interface {
	GetUserById() gin.HandlerFunc
	GetUser() gin.HandlerFunc
}

type userController struct {
	s      api.UserService
	maker  token.Maker
	config utils.Config
}

func NewUserController(s api.UserService, maker token.Maker, config utils.Config) UserController {
	return &userController{
		s:      s,
		maker:  maker,
		config: config,
	}
}

func (u *userController) GetUserById() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request types.GetUser
		if err := ctx.ShouldBindUri(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}
		user_id, err := primitive.ObjectIDFromHex(request.ID)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		user, err := u.s.GetUserById(user_id)

		if err != nil {
			ctx.JSON(http.StatusNotFound, errorRes(err))
			return
		}

		ctx.JSON(http.StatusOK, user)
	}
}

func (u *userController) GetUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authPayload := ctx.MustGet(authPayload).(*token.Payload)
		user, err := u.s.GetUserById(authPayload.UserID)

		if err != nil {
			ctx.JSON(http.StatusNotFound, errorRes(err))
			return
		}

		ctx.JSON(http.StatusOK, user)
	}
}
