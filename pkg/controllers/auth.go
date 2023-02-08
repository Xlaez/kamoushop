package controllers

import (
	"context"
	"kamoushop/pkg/models"
	"kamoushop/pkg/services/api"
	"kamoushop/pkg/services/token"
	"kamoushop/pkg/services/types"
	"kamoushop/pkg/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthController interface {
	CreateUser() gin.HandlerFunc
	LoginUser() gin.HandlerFunc
}

type authController struct {
	s      api.AuthService
	maker  token.Maker
	config utils.Config
	t_col  mongo.Collection
}

type tokens struct {
	AccessToken  string
	RefreshToken string
}

func NewAuthController(service api.AuthService, maker token.Maker, config utils.Config, token_col mongo.Collection) AuthController {
	return &authController{
		s:      service,
		maker:  maker,
		config: config,
		t_col:  token_col,
	}
}

func (a *authController) CreateUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request types.AddUser

		if err := ctx.ShouldBindJSON(&request); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		if err := a.s.CreateUser(models.User{
			FirstName: request.FirstName,
			LastName:  request.LastName,
			Email:     request.Email,
			Password:  request.Password,
		}); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		// TODO: set redis for auth codes
		ctx.JSON(http.StatusCreated, msgRes("user created successfully!"))
	}
}

func (a *authController) LoginUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request types.Login
		if err := ctx.ShouldBindJSON(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}
		user, err := a.s.Login(request)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		token, err := generateAuthTokens(ctx, a, user.ID, a.config.AccessTokenDuration)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"tokens": token})
	}
}

func (a *authController) GetUserById() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// var request
	}
}

func generateAuthTokens(ctx context.Context, a *authController, user_id primitive.ObjectID, duration time.Duration) (*tokens, error) {
	access_token, err := a.maker.CreateToken(user_id.String(), duration)
	if err != nil {
		return &tokens{}, err
	}
	refresh_token, err := a.maker.CreateToken(user_id.String(), 6000*time.Second)

	if err != nil {
		return &tokens{}, err
	}

	// TODO: Fix this repetitoion issue with insertmany
	_, err = a.t_col.InsertOne(ctx, models.Token{
		Token:     access_token,
		UserID:    user_id,
		Type:      "access",
		ExpiresAT: duration,
	})

	if err != nil {
		return &tokens{}, err
	}

	_, err = a.t_col.InsertOne(ctx, models.Token{
		Token:     refresh_token,
		UserID:    user_id,
		Type:      "refresh",
		ExpiresAT: 6000 * time.Second,
	})

	if err != nil {
		return &tokens{}, err
	}
	return &tokens{
		AccessToken:  access_token,
		RefreshToken: refresh_token,
	}, nil
}

func verifyToken(ctx context.Context, a *authController, token string) (models.Token, error) {
	var token_doc models.Token

	// TODO: check if not blacklisted
	filter := bson.D{{Key: "token", Value: token}}
	if err := a.t_col.FindOne(ctx, filter).Decode(&token_doc); err != nil {
		return models.Token{}, err
	}

	return token_doc, nil
}

func errorRes(err error) gin.H {
	return gin.H{"error: ": err.Error()}
}

func msgRes(msg string) gin.H {
	return gin.H{"messgae: ": msg}
}
