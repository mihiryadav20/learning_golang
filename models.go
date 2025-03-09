package main

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Team represents a sports team that can register on the platform
type Team struct {
	ID              string    `json:"id" gorm:"primaryKey"`
	Email           string    `json:"email" gorm:"uniqueIndex;not null"`
	Password        string    `json:"-" gorm:"not null"` // Password is not included in JSON responses
	TeamName        string    `json:"teamName" gorm:"not null"`
	Sport           string    `json:"sport" gorm:"not null"`
	TeamDescription string    `json:"teamDescription"`
	LogoPath        string    `json:"logoPath"` // Path to the stored logo file
	CreatedAt       time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt       time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (t *Team) BeforeCreate(tx *gorm.DB) error {
	if t.ID == "" {
		t.ID = uuid.New().String()
	}
	return nil
}

// TeamRegistration is used for the registration request
type TeamRegistration struct {
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=6"`
	TeamName        string `json:"teamName" validate:"required"`
	Sport           string `json:"sport" validate:"required"`
	TeamDescription string `json:"teamDescription"`
}

// TeamResponse is used for returning team data (without sensitive info)
type TeamResponse struct {
	ID              string    `json:"id"`
	Email           string    `json:"email"`
	TeamName        string    `json:"teamName"`
	Sport           string    `json:"sport"`
	TeamDescription string    `json:"teamDescription"`
	LogoPath        string    `json:"logoPath"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

// ToResponse converts a Team model to TeamResponse
func (t *Team) ToResponse() TeamResponse {
	return TeamResponse{
		ID:              t.ID,
		Email:           t.Email,
		TeamName:        t.TeamName,
		Sport:           t.Sport,
		TeamDescription: t.TeamDescription,
		LogoPath:        t.LogoPath,
		CreatedAt:       t.CreatedAt,
		UpdatedAt:       t.UpdatedAt,
	}
}
