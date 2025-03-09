package main

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Tryout represents a sports tryout event
type Tryout struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	TeamID      string    `json:"teamId" gorm:"not null;index"`
	Title       string    `json:"title" gorm:"not null"`
	Description string    `json:"description"`
	TryoutDate  time.Time `json:"tryoutDate" gorm:"not null"`
	League      string    `json:"league"`
	Division    string    `json:"division"`
	Location    string    `json:"location" gorm:"not null"`
	FormLink    string    `json:"formLink"`
	StartDate   time.Time `json:"startDate" gorm:"not null"`
	EndDate     time.Time `json:"endDate" gorm:"not null"`
	CreatedAt   time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
	LastDate    time.Time `json:"lastDate" gorm:"not null"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (t *Tryout) BeforeCreate(tx *gorm.DB) error {
	if t.ID == "" {
		t.ID = uuid.New().String()
	}
	return nil
}

// TryoutRequest represents the request to create a tryout
type TryoutRequest struct {
	Title       string `json:"title" form:"title"`
	Description string `json:"description" form:"description"`
	League      string `json:"league" form:"league"`
	Division    string `json:"division" form:"division"`
	Location    string `json:"location" form:"location"`
	LastDate    string `json:"lastDate" form:"lastDate"`     // Format: "2006-01-02"
	TryoutDate  string `json:"tryoutDate" form:"tryoutDate"` // Format: "2006-01-02"
	StartDate   string `json:"startDate" form:"startDate"`   // Format: "2006-01-02"
	EndDate     string `json:"endDate" form:"endDate"`       // Format: "2006-01-02"

}
