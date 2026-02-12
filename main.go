package main

import (
	"log"
	"story-book/internal/app"
	"story-book/internal/config"
)

// @title Story Book API
// @version 1.0
// @description API для пользователей Story Book
// @BasePath /
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	cfg := config.Load()

	if err := app.Run(cfg); err != nil {
		log.Fatal(err)
	}
}
