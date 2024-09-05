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
	if err := database.DBConn.Find(&menus).Error; err != nil {
		log.Println("Error retrieving menus:", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Error retrieving menus")
	}

	groupedMenus := make(map[string][]models.Menu)
	if len(menus) == 0 {
		log.Println("No menus found in database")
	} else {
		for _, menu := range menus {
			groupedMenus[menu.Parent] = append(groupedMenus[menu.Parent], menu)
		}
	}

	// log.Println("Grouped Menus:", groupedMenus)
	c.Locals("menus", groupedMenus)
	return c.Next()
}
