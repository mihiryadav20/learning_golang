package handlers

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"

	"gochat/models" // Replace yourusername with your GitHub username
)

// UserRepository defines the interface for user database operations
type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByUsername(username string) (*models.User, error)
}

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	userRepo UserRepository
}

// NewUserHandler creates a new user handler
func NewUserHandler(userRepo UserRepository) *UserHandler {
	return &UserHandler{
		userRepo: userRepo,
	}
}

// RegisterRequest represents the user registration request
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Register handles user registration
func (h *UserHandler) Register(c *fiber.Ctx) error {
	// Parse request body
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request data",
		})
	}

	// Basic validation
	if req.Username == "" || req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Username, email and password are required",
		})
	}

	if len(req.Password) < 6 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Password must be at least 6 characters",
		})
	}

	// Check if username already exists
	existingUser, err := h.userRepo.GetUserByUsername(req.Username)
	if err == nil && existingUser != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Username already exists",
		})
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to process registration",
		})
	}

	// Create user object
	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
		Status:   "offline",
	}

	// Save user to database
	if err := h.userRepo.CreateUser(user); err != nil {
		log.Printf("Error creating user: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	// Return response
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id":         user.ID,
		"username":   user.Username,
		"email":      user.Email,
		"created_at": user.CreatedAt,
	})
}
