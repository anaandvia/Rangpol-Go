package middleware

import (
	"rangpol/database"
	"rangpol/models"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// Authenticate checks if the provided username and password are correct
func Authenticate(c *fiber.Ctx) error {
	var user models.User
	username := c.FormValue("username")
	password := c.FormValue("password")

	// Retrieve user from the database
	if err := database.DBConn.Where("username = ?", username).First(&user).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("Invalid credentials")
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("Invalid credentials")
	}

	// Authentication successful
	return c.SendString("Login successful")
}
