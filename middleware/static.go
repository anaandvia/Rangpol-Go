package middleware

import (
	"rangpol/database"
	"rangpol/models"

	"github.com/gofiber/fiber/v2"
)

func LantaiMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Ambil data menu dari database
		var floors []models.Lantai
		if err := database.DBConn.Order("no_lantai ASC").Find(&floors).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error retrieving floors")
		}

		// Set data menu ke context
		c.Locals("floors", floors)

		return c.Next()
	}
}
