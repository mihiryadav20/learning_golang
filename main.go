package main

import (
	"log"

	"auth/config"
	"auth/database"
	"auth/routes"
	"auth/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// Load configuration
	config.LoadConfig()

	// Initialize database
	database.InitDatabase()

	// Start token blacklist cleanup routine
	utils.StartBlacklistCleanup()

	// Create Fiber app
	app := fiber.New()

	// Middleware
	app.Use(cors.New())
	app.Use(logger.New())

	// Setup routes
	routes.SetupRoutes(app)

	// Start server
	port := config.GetEnv("PORT", "3000")
	log.Printf("Starting server on port %s...\n", port)
	log.Fatal(app.Listen(":" + port))
}
