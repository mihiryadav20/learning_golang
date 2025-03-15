package handlers

import (
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/golang-jwt/jwt/v4"

	"gochat/chat" // Replace with your GitHub username
)

// ChatHub is the global chat hub instance
var ChatHub *chat.ChatHub

// InitChatHub initializes the chat hub
func InitChatHub(userRepo interface{}) {
	// Type assertion to get the correct user repository type
	userRepoTyped, ok := userRepo.(chat.UserRepository)
	if !ok {
		log.Fatalf("Invalid user repository type passed to InitChatHub")
	}

	ChatHub = chat.NewChatHub(userRepoTyped)
	ChatHub.Run()
}

// WebSocketHandler handles WebSocket connections
func WebSocketHandler(c *websocket.Conn) {
	// Get user ID from locals
	userID, ok := c.Locals("userID").(int64)
	if !ok {
		log.Println("Missing or invalid userID in WebSocket handler")
		return
	}

	// Hand over to the chat hub
	ChatHub.HandleWebSocket(c, userID)
}

// WebSocketMiddleware authenticates WebSocket connections
func WebSocketMiddleware(c *fiber.Ctx) error {
	// Check if it's a WebSocket upgrade request
	if websocket.IsWebSocketUpgrade(c) {
		// Get the token from query parameter
		token := c.Query("token")
		if token == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "No authentication token provided")
		}

		// Parse and validate the token
		parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			// Replace with your actual secret key
			return []byte("your-secret-key"), nil
		})

		if err != nil || !parsedToken.Valid {
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid authentication token")
		}

		// Extract user ID from claims
		claims, ok := parsedToken.Claims.(jwt.MapClaims)
		if !ok {
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid token claims")
		}

		// Convert user_id to int64
		var userID int64
		switch id := claims["user_id"].(type) {
		case float64:
			userID = int64(id)
		case int64:
			userID = id
		case string:
			var err error
			userID, err = strconv.ParseInt(id, 10, 64)
			if err != nil {
				return fiber.NewError(fiber.StatusUnauthorized, "Invalid user ID in token")
			}
		default:
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid user ID in token")
		}

		// Store user ID in locals for the WebSocket handler
		c.Locals("userID", userID)

		// Allow the upgrade
		return c.Next()
	}

	return fiber.NewError(fiber.StatusUpgradeRequired, "WebSocket upgrade required")
}
