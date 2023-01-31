package server

import (
	"context"
	"fmt"
	"kamoushop/pkg/controllers"
	"kamoushop/pkg/routes"
	"kamoushop/pkg/services/api"
	"kamoushop/pkg/services/token"
	"kamoushop/pkg/utils"
	"log"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	tokenMaker      token.Maker
	auth_controller controllers.AuthController
)

func InitTokenMaker(config utils.Config) error {
	var err error
	tokenMaker, err = token.NewPasetoMaker(config.TokenKey)

	if err != nil {
		return fmt.Errorf("cannot create the token maker: %w", err)
	}
	return nil
}

func InitCols(client *mongo.Client, config utils.Config, ctx context.Context) (*controllers.AuthController, token.Maker) {
	users_col := client.Database(config.DbName).Collection(config.UserCol)
	// products_col := client.Database(config.DbName).Collection(config.ProductCol)

	user_Service := api.NewAuthService(users_col, ctx)
	auth_controller = controllers.NewAuthController(user_Service, tokenMaker, config)
	return &auth_controller, tokenMaker
}

func Run() *gin.Engine {
	config, err := utils.LoadConfig(".")

	if err != nil {
		log.Fatal("cannot load env", err)
	}

	ctx := context.TODO()
	InitTokenMaker(config)

	mongoConn := options.Client().ApplyURI(config.MongoUri)
	mongoClient, err := mongo.Connect(ctx, mongoConn)

	if err != nil {
		log.Panic((err.Error()))
	}

	if err := mongoClient.Ping(ctx, readpref.Primary()); err != nil {
		log.Panic((err.Error()))
	}

	fmt.Println("MongoDB connection succesful!")

	auth_col, token_maker := InitCols(mongoClient, config, ctx)
	server := gin.Default()
	routes.AuthRoutes(server, *auth_col, token_maker)
	return server
}
