package chat

import (
	"gochat/models" // Replace with your GitHub username
)

// UserRepository defines the interface for the user repository needed by the chat hub
type UserRepository interface {
	GetUserByID(id int64) (*models.User, error)
	UpdateUserStatus(id int64, status string) error
}
