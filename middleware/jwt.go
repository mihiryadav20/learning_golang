package middleware

import (
	"auth/config"
	"auth/utils"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

// JWTProtected returns JWT middleware
func JWTProtected() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:   config.GetJWTSecret(),
		ErrorHandler: jwtError,
		SuccessHandler: func(c *fiber.Ctx) error {
			// Get the token string
			authHeader := c.Get("Authorization")
			tokenString := authHeader[7:] // Skip "Bearer "

			// Check if token is blacklisted
			if utils.IsTokenBlacklisted(tokenString) {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "Token has been invalidated, please login again",
				})
			}

			return c.Next()
		},
	})
}

// JWT error handler
func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing or malformed JWT",
		})
	}
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error": "Invalid or expired JWT",
	})
}
