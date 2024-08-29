package controller

import (
	"rangpol/middleware"

	"github.com/gofiber/fiber/v2"
)

// HomeController handles the request for the home page.
func LoginFormController(c *fiber.Ctx) error {
	// Ambil sesi pengguna
	sess, err := middleware.GetSessionStore().Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error retrieving session")
	}

	// Ambil pesan flash error jika ada
	flashError := sess.Get("flash_error")
	sess.Delete("flash_error")
	if err := sess.Save(); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error saving session")
	}

	// Render template login dengan pesan flash error
	return c.Render("login", fiber.Map{
		"flash_error": flashError,
	})
}

func RegisterFormController(c *fiber.Ctx) error {
	// Render the template using Fiber's render method
	return c.Render("register", nil)
}
