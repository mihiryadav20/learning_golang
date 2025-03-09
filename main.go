package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Initialize database
	InitDatabase()

	// Create a new Fiber app
	app := fiber.New(fiber.Config{
		AppName: "Tryouts API",
	})

	// Use middleware
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New())

	// Setup a simple test route at the root level
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "pong",
		})
	})

	// Setup routes
	setupRoutes(app)

	// Start the server
	log.Println("Server starting on :3000")
	log.Fatal(app.Listen(":3000"))
}

// setupRoutes configures all application routes
func setupRoutes(app *fiber.App) {
	// Root route for testing
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "API is running",
		})
	})

	// API group
	api := app.Group("/api")

	// Simple test route
	api.Get("/test", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "Test route is working",
		})
	})

	// Health check
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})

	// Public routes
	// Team routes
	teams := api.Group("/teams")

	// Register a new team
	teams.Post("/register", registerTeam)

	// Login
	teams.Post("/login", LoginTeam)

	// Get team by ID (public)
	teams.Get("/:id", getTeam)

	// Add a simple logout route to teams (for testing)
	teams.Post("/simplelogout", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "Simple logout successful",
		})
	})

	// Protected routes (require authentication)
	// Auth group
	auth := api.Group("/auth")

	// Add a direct test route before middleware
	auth.Get("/test", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "Auth route is accessible without middleware",
		})
	})

	// Apply authentication middleware
	auth.Use(AuthMiddleware())

	// Get current team
	auth.Get("/me", GetCurrentTeam)

	// Logout - implement directly for testing
	auth.Post("/logout", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "Logged out successfully",
		})
	})

	// Tryout routes (protected)
	tryoutsAuth := auth.Group("/tryouts")
	tryoutsAuth.Post("/", CreateTryout)
	tryoutsAuth.Get("/my", GetTryoutsByTeam)

	// Public tryout routes
	tryouts := api.Group("/tryouts")             // Get all tryouts with pagination
	tryouts.Get("/upcoming", GetUpcomingTryouts) // Get upcoming tryouts
	tryouts.Get("/search", SearchTryouts)        // Search tryouts
	tryouts.Get("/:id", GetTryout)
	// Update team (protected)
	// To be implemented

	// Serve static files from the uploads directory
	app.Static("/uploads", "./uploads")
}
