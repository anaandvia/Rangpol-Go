package controller

import (
	"fmt"
	"rangpol/database"
	_ "rangpol/helper"
	"rangpol/middleware"

	"rangpol/models"

	"github.com/gofiber/fiber/v2"
)

// AdminPage handles the /admin route
func DataAdminController(c *fiber.Ctx) error {
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
	if err := database.DBConn.Find(&users).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error retrieving floors")
	}

	// Set header untuk menonaktifkan caching pada halaman ini
	c.Set("Cache-Control", "no-store")
	c.Set("Pragma", "no-cache")
	c.Set("Expires", "0")

	// Render halaman dengan data yang diperlukan
	fmt.Println("menus : ", menus)

	// add := helper.add()

	return c.Render("dataadmin", fiber.Map{
		"isIndex":       1,
		"Dashboard":     "Data Admin",
		"flash_error":   flashError,
		"flash_success": flashSuccess,
		"Title":         "Data Admin",
		"menus":         menus,
		"Name":          userName,
		"User":          users,
	})
}
