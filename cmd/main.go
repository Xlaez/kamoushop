package main

import (
	_ "kamoushop/docs"
	"kamoushop/pkg/server"
	"kamoushop/pkg/utils"

	"log"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title KamouShop API
// @version 1.0
// @description This is a mini-online store that provides the basic features which one should

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:4141
// @BasePath  /v1

// @securityDefinitions.basic  BasicAuth
func main() {
	config, err := utils.LoadConfig(".")

	if err != nil {
		log.Fatal("cannot not load config", err)
	}

	s := server.Run()
	s.GET("/v1/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Fatal(s.Run(":" + config.Port))
}
