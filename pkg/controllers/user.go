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
	"golang.org/x/crypto/bcrypt"
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
	GetAllUsers() gin.HandlerFunc
	QueryBrands() gin.HandlerFunc
	DeleteUser() gin.HandlerFunc
	StarUserShop() gin.HandlerFunc
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

// GetUserById godoc
// @Summary Get a user by _id
// @Tags user
// @Accept json
// @Produce json
// @Param types.GetUser query types.GetUser true "validation code"
// @Success 200 {string} types.User
// @Router		/user/by-id/:id	[get]
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

// GetUserById godoc
// @Summary Get a user by _id
// @Tags user
// @Accept json
// @Produce json
// @Param types.GetUser query types.GetUser true "validation code"
// @Success 200 {string} types.User
// @Router		/user/by-id/:id	[get]
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

// ChangePassword godoc
// @Summary Change user's password
// @Tags user
// @Accept json
// @Produce json
// @Param types.ChangePassword body types.ChangePassword true "change user password"
// @Success 200 {string} msgRes
// @Router		/user/update/password	[patch]
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

// UpdateImage godoc
// @Summary Change user's image
// @Tags user
// @Accept json
// @Produce json
// @Param types.UploadImage formData types.UploadImage true "change user image"
// @Success 200 {string} msgRes
// @Router		/user/update/image	[patch]
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

// UpdateProfile godoc
// @Summary Change user's profile
// @Tags user
// @Accept json
// @Produce json
// @Param types.UpdateProfile body types.UpdateProfile true "change user profile"
// @Success 200 {string} msgRes
// @Router		/user/update/profile	[patch]
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

// UpdateBrandName godoc
// @Summary Update user's brand name
// @Tags user
// @Accept json
// @Produce json
// @Param types.UpdateBrandName body types.UpdateBrandName true "update user's brand name"
// @Success 200 {string} msgRes
// @Router		/user/update/brand-name	[patch]
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

// GetAllUsers godoc
// @Summary Get all the users from the database
// @Tags user
// @Accept json
// @Produce json
// @Param types.GetUsers query types.GetUsers true "get all users from database"
// @Success 200 {string} msgRes
// @Router		/user	[get]
func (u *userController) GetAllUsers() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request types.GetUsers

		if err := ctx.ShouldBindQuery(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		users, totalDocs, err := u.s.GetAllUsers(request.Limit, request.Page)

		if err != nil {
			if err == mongo.ErrNoDocuments {
				ctx.JSON(http.StatusNotFound, errorRes(errors.New("resource not found")))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"users": users, "totalDocuments": totalDocs})
	}
}

// QueryBrands godoc
// @Summary Query brands from database
// @Tags user
// @Accept json
// @Produce json
// @Param types.QueryBrands query types.QueryBrands true "get all brands from database"
// @Success 200 {string} msgRes
// @Router		/user/brands	[get]
func (u *userController) QueryBrands() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request types.QueryBrands

		if err := ctx.ShouldBindQuery(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		brands, totalDocs, err := u.s.QueryBrands(request.Keyword, request.Limit, request.Page)

		if err != nil {
			if err == mongo.ErrNoDocuments {
				ctx.JSON(http.StatusNotFound, errorRes(errors.New("resource not found")))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"brands": brands, "totalDocuments": totalDocs})
	}
}

// DeleteUser godoc
// @Summary Delete a user from database
// @Tags user
// @Accept json
// @Produce json
// @Param types.DeleteUser query types.DeleteUser true "delete user"
// @Success 200 {string} msgRes
// @Router		/user/:password	[delete]
func (u *userController) DeleteUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request types.DeleteUser
		if err := ctx.ShouldBindUri(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		payload := ctx.MustGet(authPayload).(*token.Payload)

		user, err := u.s.FindOne(bson.D{primitive.E{Key: "_id", Value: payload.UserID}})

		if err != nil {
			ctx.JSON(http.StatusExpectationFailed, errorRes(err))
			return
		}

		if err = password.ComparePassword(request.Password, user.Password); err != nil {
			if err == bcrypt.ErrMismatchedHashAndPassword {
				ctx.JSON(http.StatusBadRequest, errorRes(errors.New("password doesn't match")))
				return
			}
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		if err = u.s.DeleteUser(payload.UserID); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		ctx.JSON(http.StatusOK, msgRes("deleted"))
	}
}

// StarUserShop godoc
// @Summary Star a user's shop
// @Tags user
// @Accept json
// @Produce json
// @Param types.StarShop query types.StarShop true "user id"
// @Success 200 {string} msgRes
// @Router		/user/star/:id	[patch]
func (u *userController) StarUserShop() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request types.StarShop
		if err := ctx.ShouldBindUri(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		shop_id, err := primitive.ObjectIDFromHex(request.ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		payload := ctx.MustGet(authPayload).(*token.Payload)
		filter := bson.D{primitive.E{Key: "_id", Value: shop_id}}
		updateObj := bson.D{{Key: "$inc", Value: bson.D{{Key: "stars", Value: 1}}}, {Key: "$addToSet", Value: bson.D{{Key: "starredBy", Value: payload.UserID}}}}

		if err = u.s.UpdateUser(filter, updateObj); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		ctx.JSON(http.StatusOK, msgRes("updated"))
	}
}
