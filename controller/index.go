package controller

import (
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

	// Ambil pesan flash error jika ada
	flashError := sess.Get("flash_error")
	sess.Delete("flash_error")
	if err := sess.Save(); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error saving session")
	}

	// Ambil data lantai dari database
	// var floors []models.Lantai
	// if err := database.DBConn.Order("no_lantai ASC").Find(&floors).Error; err != nil {
	// 	return c.Status(fiber.StatusInternalServerError).SendString("Error retrieving floors")
	// }
	idLantai := c.Query("id")

	floors := c.Locals("floors").([]models.Lantai)

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

	return c.Render("index", fiber.Map{
		"flash_error": flashError,
		"Title":       "Home Page",
		"Floors":      floors,
		"Rooms":       rooms,
		"isIndex":     1,
		"Dashboard":   dashboard,
	})
}
