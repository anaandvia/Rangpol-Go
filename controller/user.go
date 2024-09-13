package controller

import (
	"rangpol/database"
	"rangpol/middleware"
	"rangpol/models"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// RegisterUser handles user registration
func RegisterController(c *fiber.Ctx) error {
	// Retrieve form values
	nim := c.FormValue("nim")
	nameUser := c.FormValue("name_user")
	email := c.FormValue("email")
	username := c.FormValue("username")
	password := c.FormValue("password")
	fotoUser := c.FormValue("foto_user")
	code := c.FormValue("code")

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error hashing password")
	}

	// Create user instance
	user := models.User{
		Nim:         nim,
		Name_user:   nameUser,
		Email:       email,
		Username:    username,
		Password:    string(hashedPassword),
		Foto_user:   fotoUser,
		Verif_email: false, // Assuming default value is false
		Code:        code,
	}

	// Save user to the database
	if err := database.DBConn.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error creating user")
	}

	return c.Redirect("/login")
}

func LoginPostController(c *fiber.Ctx) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	var user models.User
	if err := database.DBConn.Where("username = ?", username).First(&user).Error; err != nil {
		sess, _ := middleware.GetSessionStore().Get(c)
		sess.Set("flash_error", "Username atau Password Tidak Sesuai")
		sess.Save()
		return c.Redirect("/login")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		sess, _ := middleware.GetSessionStore().Get(c)
		sess.Set("flash_error", "Username atau Password Tidak Sesuai")
		sess.Save()
		return c.Redirect("/login")
	}

	// Set session user_id after successful login
	sess, _ := middleware.GetSessionStore().Get(c)
	sess.Set("user_id", user.Id_user)
	sess.Set("name_user", user.Name_user)
	sess.Set("role_id", user.Level)

	sess.Save()

	return c.Redirect("/")
}

func LogoutController(c *fiber.Ctx) error {
	sess, err := middleware.GetSessionStore().Get(c)
	if err != nil {
		return err
	}

	// Destroy the session
	sess.Destroy()

	// Redirect to the login page
	return c.Redirect("/login")
}
