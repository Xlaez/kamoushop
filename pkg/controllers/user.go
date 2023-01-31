package controllers

import (
	"kamoushop/pkg/services/api"
	"kamoushop/pkg/services/token"
	"kamoushop/pkg/utils"

	"github.com/gin-gonic/gin"
)

type UserController interface {
	GetUser() gin.HandlerFunc
}

type userController struct {
	s      *api.UserService
	maker  token.Maker
	config utils.Config
}

func NewUserController(s *api.UserService, maker token.Maker, config utils.Config) UserController {
	return &userController{
		s:      s,
		maker:  maker,
		config: config,
	}
}

func (u *userController) GetUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}
