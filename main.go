package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/websocket/v2"

	"gochat/database"
	"gochat/handlers"
	"gochat/routes"
)

func main() {
	// Connect to database
	err := database.Connect("./chat.db")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Create a new Fiber app
	app := fiber.New()

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New())

	// Create user repository
	userRepo := database.NewUserRepository(database.DB)

	// Initialize chat hub - this is the critical line that was missing
	log.Println("Initializing chat hub...")
	handlers.InitChatHub(userRepo)
	log.Println("Chat hub initialized successfully")

	// Create handlers
	userHandler := handlers.NewUserHandler(userRepo)

	// Setup routes
	routes.SetupRoutes(app, userHandler)

	// Basic test route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("GoChat server is running! WebSocket endpoint: ws://localhost:8080/ws")
	})

	// Simple WebSocket test endpoint for debugging
	app.Get("/ws-test", websocket.New(func(c *websocket.Conn) {
		log.Println("Test WebSocket connected")

		// Simple echo for testing
		for {
			mt, msg, err := c.ReadMessage()
			if err != nil {
				log.Println("Read error:", err)
				break
			}
			log.Printf("Received on test socket: %s", msg)

			if err := c.WriteMessage(mt, msg); err != nil {
				log.Println("Write error:", err)
				break
			}
		}
	}))

	// Start server
	log.Println("Starting server on :8080")
	if err := app.Listen(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
