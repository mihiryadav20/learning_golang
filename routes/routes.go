package routes

import (
	"auth/handlers"
	"auth/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes configures the API routes
func SetupRoutes(app *fiber.App) {
	// Root route for testing
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("JWT Auth Server is running")
	})

	// Simple debug route (no group)
	app.Get("/debug-test", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "Debug route is working"})
	})

	// Auth routes
	auth := app.Group("/auth")
	auth.Post("/register", handlers.Register)
	auth.Post("/login", handlers.Login)

	// API routes with JWT middleware
	api := app.Group("/api")
	api.Use(middleware.JWTProtected())
	api.Get("/profile", handlers.GetProfile)
	api.Post("/logout", handlers.Logout)

	// Debug routes - try adding them after other routes
	debug := app.Group("/debug")
	debug.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	debug.Get("/routes", func(c *fiber.Ctx) error {
		routeCount := len(app.GetRoutes())
		return c.JSON(fiber.Map{
			"message":      "Debug routes endpoint reached",
			"routes_count": routeCount,
		})
	})
}
