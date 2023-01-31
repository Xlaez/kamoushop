package controllers

import (
	"kamoushop/pkg/services/api"
	"kamoushop/pkg/services/token"
	"kamoushop/pkg/utils"

	"github.com/gin-gonic/gin"
)

type AuthController interface {
	CreateUser() gin.HandlerFunc
}

type authController struct {
	s      api.AuthService
	maker  token.Maker
	config utils.Config
}

func NewAuthController(service api.AuthService, maker token.Maker, config utils.Config) AuthController {
	return &authController{
		s:      service,
		maker:  maker,
		config: config,
	}
}

func (c *authController) CreateUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}
