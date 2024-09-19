package controller

import (
	"fmt"
	"log"
	"rangpol/database"
	_ "rangpol/helper"
	"rangpol/middleware"
	"strconv"

	"rangpol/models"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// AdminPage handles the /admin route
func DataUserController(c *fiber.Ctx) error {
	// if err := middleware.RequireUserLevel(c, 1); err != nil {
	// 	return err
	// }

	// Serve the admin page
	// Ambil sesi pengguna
	sess, err := middleware.GetSessionStore().Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error retrieving session")
	}

	// Periksa apakah user_id ada di sesi
	userID := sess.Get("user_id")
	if userID == nil {
		return c.Redirect("/login")
	}

	userName := sess.Get("name_user")

	// Ambil pesan flash error
	flashError := sess.Get("flash_error")

	// Ambil pesan flash success
	flashSuccess := sess.Get("flash_success")

	// Hapus kedua pesan flash setelah diambil
	sess.Delete("flash_error")
	sess.Delete("flash_success")

	// Simpan session setelah semua operasi selesai
	if err := sess.Save(); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error saving session")
	}

	menus := c.Locals("menus").([]models.Menu)

	var users []models.User
	if err := database.DBConn.Where("dlt = ?", 0).Find(&users).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error retrieving floors")
	}

	// Set header untuk menonaktifkan caching pada halaman ini
	c.Set("Cache-Control", "no-store")
	c.Set("Pragma", "no-cache")
	c.Set("Expires", "0")

	// Render halaman dengan data yang diperlukan
	fmt.Println("menus : ", menus)

	// add := helper.add()

	return c.Render("datauser", fiber.Map{
		"isIndex":       1,
		"Dashboard":     "Data User",
		"flash_error":   flashError,
		"flash_success": flashSuccess,
		"Title":         "Data User",
		"menus":         menus,
		"Name":          userName,
		"User":          users,
	})
}

func CreateUserController(c *fiber.Ctx) error {
	sess, err := middleware.GetSessionStore().Get(c)
	if err != nil {
		log.Println("Error getting session:", err)
		return err
	}

	// Retrieve form values
	nim := c.FormValue("nim")
	nameUser := c.FormValue("name_user")
	email := c.FormValue("email")
	username := c.FormValue("nim")
	password := c.FormValue("nim") // Changed from "username" to "password"
	fotoUser := c.FormValue("foto_user")
	code := c.FormValue("code")
	Level, _ := strconv.Atoi(c.FormValue("level"))

	// Check if nim already exists
	var existingUser models.User
	if err := database.DBConn.Where("nim = ? AND dlt = ? ", nim, 0).First(&existingUser).Error; err == nil {
		// Nim already exists
		log.Printf("Nim '%s' already exists", nim)
		sess.Set("flash_error", "Nim already exists")
		if err := sess.Save(); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error saving session")
		}
		return c.Redirect("/admin/datauser")
	}

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
		Level:       uint(Level),
	}

	// Save user to the database
	if err := database.DBConn.Create(&user).Error; err != nil {
		sess.Set("flash_error", "Error creating user")
		if err := sess.Save(); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error saving session")
		}
		return c.Status(fiber.StatusInternalServerError).SendString("Error creating user")
	}

	sess.Set("flash_success", "User created successfully")

	// Save session
	if err := sess.Save(); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error saving session")
	}

	return c.Redirect("/admin/datauser")
}

func DeleteUserController(c *fiber.Ctx) error {
	sess, err := middleware.GetSessionStore().Get(c)
	if err != nil {
		log.Println("Error getting session:", err)
		return err
	}

	// Retrieve user ID from the URL path
	userID := c.Params("id")

	// Find the user by ID
	var user models.User
	if err := database.DBConn.First(&user, userID).Error; err != nil {
		sess.Set("flash_error", "User not found")
		if err := sess.Save(); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error saving session")
		}
		return c.Redirect("/admin/datauser")
	}

	// Mark user as deleted (set dlt to 1)
	user.Dlt = 1
	if err := database.DBConn.Save(&user).Error; err != nil {
		sess.Set("flash_error", "Error updating user")
		if err := sess.Save(); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error saving session")
		}
		return c.Status(fiber.StatusInternalServerError).SendString("Error updating user")
	}

	sess.Set("flash_success", "User marked as deleted successfully")

	// Save session
	if err := sess.Save(); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error saving session")
	}

	return c.Redirect("/admin/datauser")
}

func UpdateUserController(c *fiber.Ctx) error {
	sess, err := middleware.GetSessionStore().Get(c)
	if err != nil {
		log.Println("Error getting session:", err)
		return err
	}

	// Retrieve user ID from the URL path
	userID := c.Params("id")

	// Retrieve form values
	nim := c.FormValue("nim")
	nameUser := c.FormValue("name_user")
	email := c.FormValue("email")
	username := c.FormValue("username")
	fotoUser := c.FormValue("foto_user")
	code := c.FormValue("code")
	level, _ := strconv.Atoi(c.FormValue("level"))
	dlt, _ := strconv.Atoi(c.FormValue("dlt")) // For managing deletion status

	var existingUser models.User
	if err := database.DBConn.Where("nim = ? AND id_user != ? AND dlt = ?", nim, userID, 0).First(&existingUser).Error; err == nil {
		// NIM is already used by another user
		sess.Set("flash_error", "NIM already in use")
		if err := sess.Save(); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error saving session")
		}
		return c.Redirect("/admin/datauser")
	}

	// Find the user by ID
	var user models.User
	if err := database.DBConn.First(&user, userID).Error; err != nil {
		sess.Set("flash_error", "User not found")
		if err := sess.Save(); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error saving session")
		}
		return c.Redirect("/admin/datauser")
	}

	// Update user details
	user.Nim = nim
	user.Name_user = nameUser
	user.Email = email
	user.Username = username
	user.Foto_user = fotoUser
	user.Code = code
	user.Level = uint(level)
	user.Dlt = dlt // Update deletion status if applicable

	// Save updated user to the database
	if err := database.DBConn.Save(&user).Error; err != nil {
		sess.Set("flash_error", "Error updating user")
		if err := sess.Save(); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error saving session")
		}
		return c.Status(fiber.StatusInternalServerError).SendString("Error updating user")
	}

	sess.Set("flash_success", "User updated successfully")

	// Save session
	if err := sess.Save(); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error saving session")
	}

	return c.Redirect("/admin/datauser")
}
