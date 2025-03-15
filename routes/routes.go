package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"

	"gochat/handlers" // Replace with your GitHub username
)

// SetupRoutes configures all application routes
func SetupRoutes(app *fiber.App, userHandler *handlers.UserHandler) {
	// API group
	api := app.Group("/api")

	// Auth routes
	auth := api.Group("/auth")
	auth.Post("/register", userHandler.Register)
	auth.Post("/login", userHandler.Login)

	// User routes (to be implemented later)
	users := api.Group("/users")
	users.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "List users endpoint - to be implemented"})
	})

	// WebSocket configuration
	// First add the middleware for authentication
	app.Use("/ws", handlers.WebSocketMiddleware)

	// Then add the WebSocket handler
	app.Get("/ws", websocket.New(handlers.WebSocketHandler))
}
