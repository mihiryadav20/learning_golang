package database

import (
	"log"

	"auth/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DB is the database instance
var DB *gorm.DB

// InitDatabase initializes the database connection
func InitDatabase() {
	var err error
	log.Println("Opening SQLite database connection...")
	DB, err = gorm.Open(sqlite.Open("auth.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	log.Println("Connected to database successfully")

	// Auto migrate the schema
	log.Println("Running database migrations...")
	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	log.Println("Database migrations completed successfully")
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}
