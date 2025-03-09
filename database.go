package main

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDatabase initializes the database connection
func InitDatabase() {
	var err error

	// For development, we're using SQLite
	// In production, you would use PostgreSQL, MySQL, etc.
	DB, err = gorm.Open(sqlite.Open("tryouts.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto migrate the schema
	err = DB.AutoMigrate(&Team{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	//Tryouts table

	err = DB.AutoMigrate(&Team{}, &Tryout{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	log.Println("Database connected and migrated successfully")
}
