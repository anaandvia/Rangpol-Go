package middleware

import (
	"log"
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

func GetMenu(c *fiber.Ctx) error {
	var menus []models.Menu
	if err := database.DBConn.Preload("Children").Order("urutan ASC").Find(&menus).Order("Urutan asc").Error; err != nil {
		log.Println("Error retrieving menus:", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Error retrieving menus")
	}

	// log.Println("Grouped Menus:", groupedMenus)
	c.Locals("menus", menus)
	return c.Next()
}
