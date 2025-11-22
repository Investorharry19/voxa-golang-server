package main

// @title Voxa API
// @version 1.0
// @description Voxa Golang Server API

// @host localhost:3000
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter your token with the "Bearer " prefix, e.g. "Bearer eyJhbGciOi..."

import (
	"fmt"

	_ "github.com/Investorharry19/voxa-golang-server/docs"

	swagger "github.com/swaggo/fiber-swagger"

	"github.com/Investorharry19/voxa-golang-server/config"
	"github.com/Investorharry19/voxa-golang-server/database"
	"github.com/Investorharry19/voxa-golang-server/routers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	fmt.Println("Hello Voxa")

	app := fiber.New()
	app.Use(logger.New())

	app.Get("/swagger/*", swagger.WrapHandler)
	routers.UserRouter(app)
	routers.MessageRouter(app)

	config.InitCloudinary()
	database.ConnectMongoDB()
	app.Listen(":3000")
}
