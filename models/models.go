package models

import (
	"time"
)

// User represents a chat application user
type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`      // Password is never returned in JSON
	Status    string    `json:"status"` // online, offline, away
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Message represents a chat message
type Message struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// Room represents a chat room
type Room struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// RoomMember represents a user membership in a room
type RoomMember struct {
	UserID   int64     `json:"user_id"`
	RoomID   int64     `json:"room_id"`
	JoinedAt time.Time `json:"joined_at"`
}
