package main

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// registerTeam handles team registration
func registerTeam(c *fiber.Ctx) error {
	// Parse request body
	registration := new(TeamRegistration)
	if err := c.BodyParser(registration); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	// Validate required fields
	if registration.Email == "" || registration.Password == "" ||
		registration.TeamName == "" || registration.Sport == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Required fields: email, password, teamName, sport",
		})
	}

	// Validate email format
	if !isValidEmail(registration.Email) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid email format",
		})
	}

	// Check if email already exists
	var existingTeam Team
	result := DB.Where("email = ?", strings.ToLower(registration.Email)).First(&existingTeam)
	if result.Error == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"status":  "error",
			"message": "Email already registered",
		})
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Database error",
		})
	}

	// Check password strength
	if len(registration.Password) < 6 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Password must be at least 6 characters long",
		})
	}

	// Hash password
	hashedPassword, err := hashPassword(registration.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Could not process registration",
		})
	}

	// Create new team
	newTeam := Team{
		Email:           strings.ToLower(registration.Email),
		Password:        hashedPassword,
		TeamName:        registration.TeamName,
		Sport:           registration.Sport,
		TeamDescription: registration.TeamDescription,
	}

	// Save to database
	if err := DB.Create(&newTeam).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to register team",
			"error":   err.Error(),
		})
	}

	// Return success response
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Team registered successfully",
		"data":    newTeam.ToResponse(),
	})
}

// getTeam returns a team by ID
func getTeam(c *fiber.Ctx) error {
	id := c.Params("id")

	// Find team in database
	var team Team
	result := DB.First(&team, "id = ?", id)

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

// Helper Functions

// hashPassword hashes a password using bcrypt
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// isValidEmail performs basic email validation
func isValidEmail(email string) bool {
	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return false
	}

	parts := strings.Split(email, "@")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return false
	}

	return true
}

// verifyPassword checks if the provided password matches the hash
func verifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
