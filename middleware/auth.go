// middleware/auth.go
package middleware

import (
	"rangpol/database"
	"rangpol/models"

	"github.com/gofiber/fiber/v2"
)

func RedirectIfAuthenticated(c *fiber.Ctx) error {
	sess, err := GetSessionStore().Get(c)
	if err != nil {
		return err
	}

	// Check if the user is authenticated
	if sess.Get("user_id") != nil {
		return c.Redirect("/")
	}

	return c.Next()
}

func RequireUserLevel(c *fiber.Ctx, requiredLevel uint) error {
	sess, err := GetSessionStore().Get(c)
	if err != nil {
		return err
	}

	userID := sess.Get("user_id")
	if userID == nil {
		sess.Set("flash_error", "Unauthorized access. Please log in.")
		if err := sess.Save(); err != nil {
			return err
		}
		return c.Redirect("/login")
	}

	var user models.User
	if err := database.DBConn.First(&user, userID).Error; err != nil {
		sess.Set("flash_error", "Unauthorized access. Please log in.")
		if err := sess.Save(); err != nil {
			return err
		}
		return c.Redirect("/login")
	}

	if user.Level != requiredLevel {
		sess.Set("flash_error", "Access denied. You do not have the required permissions.")
		if err := sess.Save(); err != nil {
			return err
		}
		return c.Redirect("/")
	}

	return nil
}
