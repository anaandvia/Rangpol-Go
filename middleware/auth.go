// middleware/auth.go
package middleware

import (
	"log"
	"rangpol/database"
	"rangpol/models"

	"github.com/gofiber/fiber/v2"
)

func RedirectIfAuthenticated(c *fiber.Ctx) error {
	sess, err := GetSessionStore().Get(c)
	if err != nil {
		return err
	}

	userID := sess.Get("user_id")
	roleID := sess.Get("role_id")

	// Jika user sudah login, redirect ke halaman sesuai dengan role
	if c.Path() == "/login" {
		log.Println("Checking login path...")

		// If user is authenticated, redirect to the appropriate page
		if userID != nil {
			log.Println("userID is not nil")

			// Assert roleID as uint
			roleIDUint, ok := roleID.(uint)
			if !ok {
				log.Println("roleID is not of type uint")
				return c.Redirect("/")
			}

			// Compare the uint value
			if roleIDUint == 2 {
				log.Println("Redirecting to /admin")
				return c.Redirect("/admin")
			}

			return c.Redirect("/")
		}

		// Continue to login page if user is not authenticated
		log.Println("user is not authenticated")
		return c.Next()
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
