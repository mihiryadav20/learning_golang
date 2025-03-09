package main

import (
	"github.com/gofiber/fiber/v2"
)

// ListAllRoutes returns all registered routes
func ListAllRoutes(app *fiber.App) fiber.Handler {
	return func(c *fiber.Ctx) error {
		routes := make([]map[string]string, 0)

		for _, r := range app.GetRoutes() {
			routes = append(routes, map[string]string{
				"method": r.Method,
				"path":   r.Path,
			})
		}

		return c.JSON(fiber.Map{
			"status": "success",
			"count":  len(routes),
			"routes": routes,
		})
	}
}

// CheckDatabase tests the database connection
func CheckDatabase(c *fiber.Ctx) error {
	var count int64
	result := DB.Model(&Team{}).Count(&count)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Database error",
			"error":   result.Error.Error(),
		})
	}

	// Get some basic information about the database
	var tables []string
	DB.Raw("SELECT name FROM sqlite_master WHERE type='table'").Pluck("name", &tables)

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Database connection successful",
		"data": fiber.Map{
			"teamsCount": count,
			"tables":     tables,
			"dialect":    DB.Dialector.Name(),
		},
	})
}

// GetServerInfo returns information about the server
func GetServerInfo(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"version":     "0.1.0",
			"environment": "development",
			"jwtExpiry":   tokenExpiration.String(),
		},
	})
}
