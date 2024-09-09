package controller

import (
	"fmt"
	"rangpol/database"
	"rangpol/middleware"
	"rangpol/models"

	"github.com/gofiber/fiber/v2"
)

// HomeController handles the request for the home page.
func HomeController(c *fiber.Ctx) error {
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

	idLantai := c.Query("id")

	floors := c.Locals("floors").([]models.Lantai)
	menus := c.Locals("menus").([]models.Menu)
	dashboard := "Dashboard"
	// Ambil data ruangan dari database
	var rooms []models.Room
	if idLantai != "" {
		if err := database.DBConn.Where("lantai = ?", idLantai).Find(&rooms).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error retrieving rooms")
		}
		dashboard = "Lantai " + idLantai
	} else {
		// Jika tidak ada ID lantai, ambil semua ruangan
		if err := database.DBConn.Find(&rooms).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error retrieving rooms")
		}
	}

	// Set header untuk menonaktifkan caching pada halaman ini
	c.Set("Cache-Control", "no-store")
	c.Set("Pragma", "no-cache")
	c.Set("Expires", "0")

	// Render halaman dengan data yang diperlukan
	fmt.Println("menus : ", menus)

	return c.Render("index", fiber.Map{
		"flash_error":   flashError,
		"flash_success": flashSuccess,
		"Title":         "Home Page",
		"Floors":        floors,
		"Rooms":         rooms,
		"isIndex":       1,
		"Dashboard":     dashboard,
		"menus":         menus,
		"Name":          userName,
	})
}
