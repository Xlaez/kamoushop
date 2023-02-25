package controllers

import (
	"context"
	"errors"
	"kamoushop/pkg/models"
	"kamoushop/pkg/services/api"
	"kamoushop/pkg/services/token"
	"kamoushop/pkg/services/types"
	"kamoushop/pkg/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthController interface {
	CreateUser() gin.HandlerFunc
	LoginUser() gin.HandlerFunc
	ValidateAcc() gin.HandlerFunc
}

type authController struct {
	s            api.AuthService
	maker        token.Maker
	config       utils.Config
	t_col        mongo.Collection
	redis_client *redis.Client
}

type tokens struct {
	AccessToken  string
	RefreshToken string
}

func NewAuthController(service api.AuthService, maker token.Maker, config utils.Config, token_col mongo.Collection, redis_client *redis.Client) AuthController {
	return &authController{
		s:            service,
		maker:        maker,
		config:       config,
		t_col:        token_col,
		redis_client: redis_client,
	}
}

// CreateUser godoc
// @Summary Create a new user
// @Tags auth
// @Accept json
// @Produce json
// @Param types.AddUser body types.AddUser true "user's data"
// @Success 201 {string} code
// @Router		/auth/register	[post]
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
			StarredBy: make([]primitive.ObjectID, 500),
		}); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}
		code, err := sendVerificationCode(ctx, a, request.Email)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}
		// TODO: send code via email
		ctx.JSON(http.StatusCreated, gin.H{"code": code})
	}
}

// LoginUser godoc
// @Summary Signin a user
// @Tags auth
// @Accept json
// @Produce json
// @Param types.Login body types.Login true "user's data"
// @Success 200 {string} token
// @Router		/auth/login	[post]
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

// ValidateAcc godoc
// @Summary Validate User's Account with validation code sent after registration
// @Tags user
// @Accept json
// @Produce json
// @Param types.ValidateAcc body types.ValidateAcc true "validation code"
// @Success 200 {string} message
// @Router		/auth/validate	[post]
func (a *authController) ValidateAcc() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request types.ValidateAcc

		if err := ctx.ShouldBindJSON(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		email, err := a.redis_client.Get(ctx, request.Code).Result()
		if err == redis.Nil {
			ctx.JSON(http.StatusNotFound, errorRes(errors.New("verification codehas expired, request for another")))
			return
		}

		if err = a.s.ValidateAcc(email); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		ctx.JSON(http.StatusNoContent, "")
	}
}

func sendVerificationCode(ctx context.Context, a *authController, email string) (string, error) {
	random_code := utils.RandomStr(6)
	//TODO: set expiration duration to 30 minutes
	if err := a.redis_client.Set(ctx, random_code, email, 0).Err(); err != nil {
		return "", err
	}
	return random_code, nil
}

func generateAuthTokens(ctx context.Context, a *authController, user_id primitive.ObjectID, duration time.Duration) (*tokens, error) {
	access_token, err := a.maker.CreateToken(user_id, duration)
	if err != nil {
		return &tokens{}, err
	}
	refresh_token, err := a.maker.CreateToken(user_id, 6000*time.Second)

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
