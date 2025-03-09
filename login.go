package main

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

const (
	// Change this to your own secret key
	jwtSecret = "your-secret-key-change-this-in-production"
	// Token expiration time
	tokenExpiration = 24 * time.Hour
)

// TokenClaims represents the claims in the JWT
type TokenClaims struct {
	TeamID string `json:"teamId"`
	jwt.RegisteredClaims
}

// LoginRequest represents the login request data
type LoginRequest struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

// LoginResponse represents the login response data
type LoginResponse struct {
	Token string       `json:"token"`
	Team  TeamResponse `json:"team"`
}

// LoginTeam handles the login process
func LoginTeam(c *fiber.Ctx) error {
	// Parse login request
	login := new(LoginRequest)

	// Try to parse from JSON first
	if err := c.BodyParser(login); err != nil {
		// If that fails, try form values
		login.Email = c.FormValue("email")
		login.Password = c.FormValue("password")
	}

	// Validate required fields
	if login.Email == "" || login.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Email and password are required",
		})
	}

	// Find team by email
	var team Team
	result := DB.Where("email = ?", strings.ToLower(login.Email)).First(&team)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid email or password",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Database error",
			"error":   result.Error.Error(),
		})
	}

	// Verify password
	if err := verifyPassword(team.Password, login.Password); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid email or password",
		})
	}

	// Generate JWT token
	token, err := generateToken(team.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Could not generate token",
			"error":   err.Error(),
		})
	}

	// Return token and team data
	return c.JSON(fiber.Map{
		"status": "success",
		"data": LoginResponse{
			Token: token,
			Team:  team.ToResponse(),
		},
	})
}

// LogoutTeam handles the logout process
func LogoutTeam(c *fiber.Ctx) error {
	// Since JWT is stateless, we don't need to invalidate the token on the server
	// In a real-world app, you might want to add the token to a blacklist
	// For now, we simply send a success response

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Logged out successfully",
	})
}

// generateToken creates a new JWT token for a team
func generateToken(teamID string) (string, error) {
	// Create token claims
	claims := TokenClaims{
		teamID,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExpiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "tryouts-api",
			Subject:   teamID,
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret key
	signedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// extractToken extracts the token from the Authorization header
func extractToken(c *fiber.Ctx) string {
	// Get Authorization header
	bearerToken := c.Get("Authorization")

	// Check if the header is in the correct format
	if len(bearerToken) > 7 && strings.ToUpper(bearerToken[0:7]) == "BEARER " {
		return bearerToken[7:]
	}

	return ""
}

// validateToken validates the JWT token
func validateToken(tokenString string) (*TokenClaims, error) {
	// Parse token
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	// Extract claims
	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// AuthMiddleware is a middleware to authenticate requests
func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Extract token from request
		tokenString := extractToken(c)
		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "Missing or invalid token",
			})
		}

		// Validate token
		claims, err := validateToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid token",
				"error":   err.Error(),
			})
		}

		// Add team ID to context for use in handlers
		c.Locals("teamID", claims.TeamID)

		// Continue to next handler
		return c.Next()
	}
}

// GetCurrentTeam returns the currently authenticated team
func GetCurrentTeam(c *fiber.Ctx) error {
	// Get team ID from context
	teamID := c.Locals("teamID").(string)

	// Get team from database
	var team Team
	result := DB.First(&team, "id = ?", teamID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Team not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Database error",
			"error":   result.Error.Error(),
		})
	}

	// Return team data
	return c.JSON(fiber.Map{
		"status": "success",
		"data":   team.ToResponse(),
	})
}
