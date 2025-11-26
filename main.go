package main

// @title Voxa API
// @version 1.0
// @description Voxa Golang Server API

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter your token with the "Bearer " prefix, e.g. "Bearer eyJhbGciOi..."

import (
	"fmt"
	"log"
	"os"

	_ "github.com/Investorharry19/voxa-golang-server/docs"
	"github.com/joho/godotenv"
	swagger "github.com/swaggo/fiber-swagger"

	"github.com/Investorharry19/voxa-golang-server/config"
	"github.com/Investorharry19/voxa-golang-server/database"
	"github.com/Investorharry19/voxa-golang-server/routers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"

	docs "github.com/Investorharry19/voxa-golang-server/docs"
)

func main() {
	fmt.Println("Hello Voxa")
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found — continuing with system environment variables")
	}

	// Determine port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // fallback for local development
	}

	// Set Swagger host dynamically
	host := os.Getenv("SWAGGER_HOST")
	if host == "" {
		host = "localhost:" + port
	}

	scheme := "http"
	if host != "localhost:"+port {
		scheme = "https"
	}

	docs.SwaggerInfo.Host = host
	docs.SwaggerInfo.Schemes = []string{scheme}

	// Initialize Fiber
	app := fiber.New()
	app.Use(logger.New())

	// Swagger endpoint
	app.Get("/swagger/*", func(c *fiber.Ctx) error {
		// Set dynamic host based on request
		docs.SwaggerInfo.Host = c.Hostname() // automatically detects localhost or live URL
		docs.SwaggerInfo.Schemes = []string{"http", "https"}
		return swagger.WrapHandler(c)
	})

	// Routers
	routers.UserRouter(app)
	routers.MessageRouter(app)

	// Config and DB
	config.InitCloudinary()
	database.ConnectMongoDB()

	// Start server
	log.Printf("Server running on port %s (Swagger Host: %s)\n", port, host)
	if err := app.Listen(":" + port); err != nil {
		log.Fatal(err)
	}
}
