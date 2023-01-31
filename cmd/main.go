package main

import (
	"kamoushop/pkg/server"
	"kamoushop/pkg/utils"
	"log"
)

func main() {
	config, err := utils.LoadConfig(".")

	if err != nil {
		log.Fatal("cannot not load config", err)
	}

	s := server.Run()
	log.Fatal(s.Run(":" + config.Port))
}
