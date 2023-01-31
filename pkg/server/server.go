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
	cors "github.com/rs/cors/wrapper/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	tokenMaker      token.Maker
	auth_controller controllers.AuthController
	user_controller controllers.UserController
)

func InitTokenMaker(config utils.Config) error {
	var err error
	tokenMaker, err = token.NewPasetoMaker(config.TokenKey)

	if err != nil {
		return fmt.Errorf("cannot create the token maker: %w", err)
	}
	return nil
}

func InitCols(client *mongo.Client, config utils.Config, ctx context.Context) (*controllers.AuthController, *controllers.UserController, token.Maker) {
	users_col := client.Database(config.DbName).Collection(config.UserCol)
	token_col := client.Database(config.DbName).Collection(config.TokenCol)
	// products_col := client.Database(config.DbName).Collection(config.ProductCol)

	auth_service := api.NewAuthService(users_col, ctx)
	user_service := api.NewUserService(users_col, ctx)
	auth_controller = controllers.NewAuthController(auth_service, tokenMaker, config, *token_col)
	user_controller = controllers.NewUserController(&user_service, tokenMaker, config)
	return &auth_controller, &user_controller, tokenMaker
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

	auth_col, users_col, token_maker := InitCols(mongoClient, config, ctx)
	server := gin.Default()
	server.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		Debug:            true, // remeber to off this for prod
		AllowedMethods:   []string{"POST", "GET", "PATCH", "DELETE", "PURGE", "OPTIONS"},
		AllowCredentials: true,
		MaxAge:           3,
	}))

	// defer mongoClient.Disconnect(ctx)

	routes.AuthRoutes(server, *auth_col, token_maker)
	routes.UserRoutes(server, *users_col, token_maker)
	return server
}
