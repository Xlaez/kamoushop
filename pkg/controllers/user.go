package controllers

import (
	"errors"
	"kamoushop/pkg/libs"
	"kamoushop/pkg/services/api"
	"kamoushop/pkg/services/password"
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

const (
	authPayload = "x-auth-token_payload"
)

type UserController interface {
	GetUserById() gin.HandlerFunc
	GetUser() gin.HandlerFunc
	ChangePassword() gin.HandlerFunc
	UpdateImage() gin.HandlerFunc
	UpdateProfile() gin.HandlerFunc
	UpdateBrandName() gin.HandlerFunc
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

func (u *userController) ChangePassword() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request types.ChangePassword

		if err := ctx.ShouldBindJSON(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		payload, _ := ctx.MustGet(authPayload).(*token.Payload)
		user, err := u.s.GetUserByIdWithPassword((payload.UserID))
		if err != nil {
			ctx.JSON(http.StatusNotFound, errorRes(err))
			return
		}

		if err = password.ComparePassword(request.OldPassword, user.Password); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(errors.New("password mismatch")))
			return
		}

		hashedPassword, err := password.HashPassword(request.NewPassword)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}
		filter := bson.D{{Key: "email", Value: user.Email}}
		updateObj := bson.D{{Key: "$set", Value: bson.D{{Key: "password", Value: hashedPassword}}}}
		if err = u.s.UpdateUser(filter, updateObj); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		ctx.JSON(http.StatusOK, msgRes("updated"))
	}
}

func (u *userController) UpdateImage() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request types.UploadImage
		if err := ctx.ShouldBind(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}
		id, err := primitive.ObjectIDFromHex(request.ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		secure_url, _, e := libs.UploadToCloud(ctx)
		if e != nil {
			ctx.JSON(http.StatusExpectationFailed, errorRes(err))
			return
		}

		filter := bson.D{primitive.E{Key: "_id", Value: id}}
		updateObj := bson.D{{Key: "$set", Value: bson.D{{Key: "image", Value: secure_url}}}}

		if err = u.s.UpdateUser(filter, updateObj); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		ctx.JSON(http.StatusOK, msgRes("uploaded"))
	}
}

func (u *userController) UpdateProfile() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request types.UpdateProfile
		if err := ctx.ShouldBindJSON(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		payload := ctx.MustGet(authPayload).(*token.Payload)

		filter := bson.D{primitive.E{Key: "_id", Value: payload.UserID}}
		uploadObj := bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "facebook", Value: request.Facebook},
				{Key: "instagram", Value: request.Instagram},
				{Key: "phoneNo", Value: request.PhoneNO},
				{Key: "updatedAt", Value: time.Now()}}}}

		if err := u.s.UpdateUser(filter, uploadObj); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		ctx.JSON(http.StatusOK, msgRes("updated"))
	}
}

func (u *userController) UpdateBrandName() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request types.UpdateBrandName
		if err := ctx.ShouldBindJSON(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		user, err := u.s.FindOne(bson.D{{Key: "brandName", Value: request.BrandName}})
		if err != nil && err != mongo.ErrNoDocuments {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		if user.IsVerified {
			ctx.JSON(http.StatusBadRequest, errorRes(errors.New("brand name taken")))
			return
		}
		payload := ctx.MustGet(authPayload).(*token.Payload)
		filter := bson.D{{Key: "_id", Value: payload.UserID}}
		updateObj := bson.D{{Key: "$set", Value: bson.D{{Key: "brandName", Value: request.BrandName}, {Key: "updatedAt", Value: time.Now()}}}}

		if err = u.s.UpdateUser(filter, updateObj); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		ctx.JSON(http.StatusOK, msgRes("updated"))
	}
}
