package controller

import (
	"rangpol/database"
	"rangpol/helper"
	"rangpol/middleware"
	"rangpol/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RoomDetailController(c *fiber.Ctx) error {

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

	// Ambil pesan flash error jika ada
	flashError := sess.Get("flash_error")
	sess.Delete("flash_error")
	if err := sess.Save(); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error saving session")
	}

	// Get the room ID from the query parameters
	id := c.Query("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Room ID is required")
	}

	var room models.Room
	// Fetch the room details from the database along with its associated details
	if err := database.DBConn.Preload("DetailRoom").First(&room, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Room not found")
		}
		return c.Status(fiber.StatusInternalServerError).SendString("Error retrieving room")
	}

	var peminjaman []models.Peminjaman

	if err := database.DBConn.
		Preload("User").
		Preload("DetailPeminjaman").
		Preload("Room").
		Where("id_room = ?", id).
		Find(&peminjaman).
		Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error retrieving data")
	}

	for i := range peminjaman {
		peminjaman[i].TglAcaraFormatted = peminjaman[i].TglAcara.Format("02-01-2006 jam 15:04")
		peminjaman[i].TglAkhirAcaraFormatted = peminjaman[i].TglAkhirAcara.Format("02-01-2006 jam 15:04")

		peminjaman[i].TglAcaraDay = helper.GetIndonesianDay(peminjaman[i].TglAcara.Weekday())
		peminjaman[i].TglAkhirAcaraDay = helper.GetIndonesianDay(peminjaman[i].TglAkhirAcara.Weekday())
	}

	floors := c.Locals("floors").([]models.Lantai)
	menus := c.Locals("menus").([]models.Menu)

	// Render the detail page with the room data
	return c.Render("detailroom", fiber.Map{
		"Title":       "Room Details",
		"Room":        room,
		"flash_error": flashError,
		"Peminjaman":  peminjaman,
		"Floors":      floors,
		"menus":       menus,
		"Name":        userName,
	})
}
