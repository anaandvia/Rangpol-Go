// sessionstore/session.go
package middleware

import (
	"rangpol/database"
	"rangpol/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

var Store *session.Store

func InitSessionStore() {
	Store = session.New(session.Config{
		CookieHTTPOnly: true,
		CookieSecure:   true,  // Pastikan HTTPS digunakan
		CookieSameSite: "Lax", // Atau "Strict" jika memungkinkan
		Expiration:     24 * time.Hour,
	})
}

func GetSessionStore() *session.Store {
	return Store
}

// controller/some_controller.go
func SomeSensitiveAction(c *fiber.Ctx) error {
	sess, err := GetSessionStore().Get(c)
	if err != nil {
		return err
	}

	userID := sess.Get("user_id")
	if userID == nil {
		return c.Redirect("/login")
	}

	var user models.User
	if err := database.DBConn.First(&user, userID).Error; err != nil {
		return err
	}

	// Ensure the user has the required level
	if user.Level != 1 {
		return c.Status(fiber.StatusForbidden).SendString("Access denied")
	}

	// Proceed with the action
	return c.SendString("Sensitive action performed")
}
