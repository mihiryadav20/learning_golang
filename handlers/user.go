package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

// GetProfile returns the user profile
func GetProfile(c *fiber.Ctx) error {
	// Get user from context
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	username := claims["username"].(string)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Welcome to the protected route",
		"user":    username,
	})
}
