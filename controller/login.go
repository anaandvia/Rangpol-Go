package controller

import (
	"rangpol/middleware"

	"github.com/gofiber/fiber/v2"
)

// HomeController handles the request for the home page.
func LoginFormController(c *fiber.Ctx) error {
	sess, err := middleware.GetSessionStore().Get(c)
	if err != nil {
		return err
	}

	// Get the flash error message, if any
	flashError := sess.Get("flash_error")
	sess.Delete("flash_error")
	sess.Save()

	// Pass the flash error message to the template
	return c.Render("login", fiber.Map{
		"flash_error": flashError,
	})
}

func RegisterFormController(c *fiber.Ctx) error {
	// Render the template using Fiber's render method
	return c.Render("register", nil)
}
