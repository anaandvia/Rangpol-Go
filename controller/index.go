package controller

import (
	"rangpol/database"
	"rangpol/middleware"
	"rangpol/models"

	"github.com/gofiber/fiber/v2"
)

// HomeController handles the request for the home page.
func HomeController(c *fiber.Ctx) error {
	// Render the template using Fiber's render method
	sess, err := middleware.GetSessionStore().Get(c)
	if err != nil {
		return err
	}

	userID := sess.Get("user_id")
	if userID == nil {
		return c.Redirect("/login")
	}

	// Get the flash error message, if any
	flashError := sess.Get("flash_error")
	sess.Delete("flash_error")
	sess.Save()

	var floors []models.Lantai
	if err := database.DBConn.Order("no_lantai ASC").Find(&floors).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error retrieving floors")
	}

	return c.Render("index", fiber.Map{
		"flash_error": flashError,
		"Title":       "Home Page",
		"Floors":      floors,
	})
}
