package main

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// CreateTryout handles creating a new tryout
func CreateTryout(c *fiber.Ctx) error {
	// Get team ID from context
	teamID := c.Locals("teamID").(string)

	// Parse request
	request := new(TryoutRequest)
	if err := c.BodyParser(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	// Debug print the entire request
	log.Printf("Received request: %+v", request)

	// Validate required fields
	if request.Title == "" || request.Location == "" ||
		request.TryoutDate == "" || request.StartDate == "" ||
		request.EndDate == "" || request.LastDate == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Required fields: title, location, tryoutDate, startDate, endDate, lastDate",
		})
	}

	// Parse dates
	tryoutDate, err := time.Parse("2006-01-02", request.TryoutDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid tryoutDate format. Use YYYY-MM-DD",
		})
	}

	startDate, err := time.Parse("2006-01-02", request.StartDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid startDate format. Use YYYY-MM-DD",
		})
	}

	endDate, err := time.Parse("2006-01-02", request.EndDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid endDate format. Use YYYY-MM-DD",
		})
	}

	lastDate, err := time.Parse("2006-01-02", request.LastDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid lastDate format. Use YYYY-MM-DD",
		})
	}

	// Validate dates
	if startDate.After(endDate) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Start date must be before end date",
		})
	}

	if tryoutDate.Before(time.Now()) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Tryout date must be in the future",
		})
	}

	if lastDate.After(tryoutDate) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Last date for registration must be before or on the tryout date",
		})
	}

	// Create tryout
	tryout := Tryout{
		TeamID:      teamID,
		Title:       request.Title,
		Description: request.Description,
		League:      request.League,
		Division:    request.Division,
		Location:    request.Location,
		FormLink:    request.FormLink, // Explicitly set FormLink
		TryoutDate:  tryoutDate,
		StartDate:   startDate,
		EndDate:     endDate,
		LastDate:    lastDate,
	}

	// Debug print before saving
	log.Printf("Tryout to be saved: %+v", tryout)

	// Save to database
	if err := DB.Create(&tryout).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create tryout",
			"error":   err.Error(),
		})
	}

	// Return success response
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Tryout created successfully",
		"data":    tryout,
	})
}

// GetTryoutsByTeam returns all tryouts for the authenticated team
func GetTryoutsByTeam(c *fiber.Ctx) error {
	// Get team ID from context
	teamID := c.Locals("teamID").(string)

	// Get tryouts from database
	var tryouts []Tryout
	result := DB.Where("team_id = ?", teamID).Order("created_at DESC").Find(&tryouts)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch tryouts",
			"error":   result.Error.Error(),
		})
	}

	// Return success response
	return c.JSON(fiber.Map{
		"status": "success",
		"data":   tryouts,
	})
}

// GetTryout returns a specific tryout
func GetTryout(c *fiber.Ctx) error {
	id := c.Params("id")

	// Get tryout from database
	var tryout Tryout
	result := DB.First(&tryout, "id = ?", id)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Tryout not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Database error",
			"error":   result.Error.Error(),
		})
	}

	// Return success response
	return c.JSON(fiber.Map{
		"status": "success",
		"data":   tryout,
	})
}

// UpdateTryout updates an existing tryout
func UpdateTryout(c *fiber.Ctx) error {
	// Get team ID from context
	teamID := c.Locals("teamID").(string)

	// Get tryout ID from params
	id := c.Params("id")

	// Check if tryout exists and belongs to the team
	var existingTryout Tryout
	result := DB.Where("id = ? AND team_id = ?", id, teamID).First(&existingTryout)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Tryout not found or you don't have permission to update it",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Database error",
			"error":   result.Error.Error(),
		})
	}

	// Parse request
	request := new(TryoutRequest)
	if err := c.BodyParser(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	// Update fields that are present in the request
	if request.Title != "" {
		existingTryout.Title = request.Title
	}

	if request.Description != "" {
		existingTryout.Description = request.Description
	}

	if request.League != "" {
		existingTryout.League = request.League
	}

	if request.Division != "" {
		existingTryout.Division = request.Division
	}

	if request.Location != "" {
		existingTryout.Location = request.Location
	}

	if request.TryoutDate != "" {
		tryoutDate, err := time.Parse("2006-01-02", request.TryoutDate)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid tryoutDate format. Use YYYY-MM-DD",
			})
		}
		existingTryout.TryoutDate = tryoutDate
	}

	if request.StartDate != "" {
		startDate, err := time.Parse("2006-01-02", request.StartDate)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid startDate format. Use YYYY-MM-DD",
			})
		}
		existingTryout.StartDate = startDate
	}

	if request.EndDate != "" {
		endDate, err := time.Parse("2006-01-02", request.EndDate)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid endDate format. Use YYYY-MM-DD",
			})
		}
		existingTryout.EndDate = endDate
	}

	if request.LastDate != "" {
		lastDate, err := time.Parse("2006-01-02", request.LastDate)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid lastDate format. Use YYYY-MM-DD",
			})
		}
		existingTryout.LastDate = lastDate
	}

	// Validate date relationships
	if existingTryout.StartDate.After(existingTryout.EndDate) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Start date must be before end date",
		})
	}

	if existingTryout.LastDate.After(existingTryout.TryoutDate) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Last date for registration must be before or on the tryout date",
		})
	}

	// Save to database
	if err := DB.Save(&existingTryout).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update tryout",
			"error":   err.Error(),
		})
	}

	// Return success response
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Tryout updated successfully",
		"data":    existingTryout,
	})
}

// DeleteTryout deletes a tryout
func DeleteTryout(c *fiber.Ctx) error {
	// Capture start time for request tracking
	startTime := time.Now()

	// Log the incoming request details
	log.Printf("[DeleteTryout] Start of request at %v", startTime)
	log.Printf("[DeleteTryout] Team ID from context: %v", c.Locals("teamID"))
	log.Printf("[DeleteTryout] Tryout ID from params: %s", c.Params("id"))

	// Get team ID from context
	teamID, ok := c.Locals("teamID").(string)
	if !ok {
		log.Printf("[DeleteTryout] Failed to retrieve team ID from context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Authentication failed",
		})
	}

	// Get tryout ID from params
	id := c.Params("id")

	// Check if tryout exists and belongs to the team
	var existingTryout Tryout
	result := DB.Where("id = ? AND team_id = ?", id, teamID).First(&existingTryout)

	if result.Error != nil {
		log.Printf("[DeleteTryout] Database query error: %v", result.Error)

		if result.Error == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Tryout not found or you don't have permission to delete it",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Database error",
			"error":   result.Error.Error(),
		})
	}

	// Log tryout details before deletion
	log.Printf("[DeleteTryout] Tryout to be deleted:")
	log.Printf("[DeleteTryout] ID: %s", existingTryout.ID)
	log.Printf("[DeleteTryout] Team ID: %s", existingTryout.TeamID)
	log.Printf("[DeleteTryout] Title: %s", existingTryout.Title)

	// Delete from database
	if err := DB.Delete(&existingTryout).Error; err != nil {
		log.Printf("[DeleteTryout] Deletion error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to delete tryout",
			"error":   err.Error(),
		})
	}

	// Log successful deletion
	log.Printf("[DeleteTryout] Tryout deleted successfully. Duration: %v", time.Since(startTime))

	// Return success response
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Tryout deleted successfully",
	})
}

// SearchTryouts searches for tryouts
func SearchTryouts(c *fiber.Ctx) error {
	// Parse query parameters
	league := c.Query("league")
	division := c.Query("division")
	location := c.Query("location")

	// Build query
	query := DB.Model(&Tryout{})

	// Apply filters
	if league != "" {
		query = query.Where("league = ?", league)
	}

	if division != "" {
		query = query.Where("division = ?", division)
	}

	if location != "" {
		query = query.Where("location LIKE ?", "%"+location+"%")
	}

	// Only show future tryouts by default
	if c.Query("includeAllDates") != "true" {
		query = query.Where("tryout_date >= ?", time.Now().Format("2006-01-02"))
	}

	// Order by tryout date
	query = query.Order("tryout_date ASC")

	// Execute query
	var tryouts []Tryout
	if err := query.Find(&tryouts).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to search tryouts",
			"error":   err.Error(),
		})
	}

	// Return results
	return c.JSON(fiber.Map{
		"status": "success",
		"data":   tryouts,
	})
}

// GetUpcomingTryouts returns tryouts with future dates
func GetUpcomingTryouts(c *fiber.Ctx) error {
	// Get current date
	today := time.Now().Format("2006-01-02")

	// Get tryouts with tryout date greater than or equal to today
	var tryouts []Tryout
	result := DB.Where("tryout_date >= ?", today).
		Order("tryout_date ASC").
		Limit(10).
		Find(&tryouts)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch upcoming tryouts",
			"error":   result.Error.Error(),
		})
	}

	// Return success response
	return c.JSON(fiber.Map{
		"status": "success",
		"data":   tryouts,
	})
}
