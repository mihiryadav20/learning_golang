package utils

import (
	"log"
	"sync"
	"time"

	"auth/config"
	"auth/models"

	"github.com/golang-jwt/jwt/v4"
)

// TokenBlacklist for storing invalidated tokens
type TokenBlacklist struct {
	sync.RWMutex
	Tokens map[string]time.Time
}

// Global token blacklist
var blacklist = TokenBlacklist{
	Tokens: make(map[string]time.Time),
}

// JWT claims struct
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// IsTokenBlacklisted checks if a token is blacklisted
func IsTokenBlacklisted(tokenString string) bool {
	blacklist.RLock()
	defer blacklist.RUnlock()

	_, exists := blacklist.Tokens[tokenString]
	return exists
}

// AddToBlacklist adds a token to the blacklist with expiry time
func AddToBlacklist(tokenString string, expiry time.Time) {
	blacklist.Lock()
	defer blacklist.Unlock()

	blacklist.Tokens[tokenString] = expiry
}

// CleanupBlacklist removes expired tokens from the blacklist
func CleanupBlacklist() {
	blacklist.Lock()
	defer blacklist.Unlock()

	now := time.Now()
	for token, expiry := range blacklist.Tokens {
		if now.After(expiry) {
			delete(blacklist.Tokens, token)
		}
	}
}

// StartBlacklistCleanup starts a routine to clean up expired tokens
func StartBlacklistCleanup() {
	ticker := time.NewTicker(15 * time.Minute)
	go func() {
		for range ticker.C {
			log.Println("Cleaning up token blacklist...")
			CleanupBlacklist()
		}
	}()
}

// CreateToken creates a new JWT token for a user
func CreateToken(user models.User) (string, error) {
	// Create token claims
	claims := Claims{
		user.Username,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "go-fiber-jwt-auth",
			Subject:   user.Username,
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token
	tokenString, err := token.SignedString(config.GetJWTSecret())
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
