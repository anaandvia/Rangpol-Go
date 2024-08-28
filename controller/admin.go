package controller

import (
	"rangpol/middleware"

	"github.com/gofiber/fiber/v2"
)

// AdminPage handles the /admin route
func AdminPage(c *fiber.Ctx) error {
	if err := middleware.RequireUserLevel(c, 1); err != nil {
		return err
	}

	// Serve the admin page
	return c.Render("admin", fiber.Map{
		"Title": "Admin Page",
	})
}
