package database

import (
	"log"

	"auth/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DB is the database instance
var DB *gorm.DB

// InitDatabase initializes the database connection
func InitDatabase() {
	var err error
	log.Println("Opening MySQL database connection...")
	dsn := "newuser:strongpassword@tcp(127.0.0.1:3306)/mydatabase?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
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
