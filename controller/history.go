package controller

import (
	"fmt"
	"rangpol/database"
	"rangpol/helper"
	"rangpol/middleware"
	"rangpol/models"

	"github.com/gofiber/fiber/v2"
)

func HistoryPeminjamanController(c *fiber.Ctx) error {

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
	flashSuccess := sess.Get("flash_success")

	sess.Delete("flash_error")
	sess.Delete("flash_success")

	if err := sess.Save(); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error saving session")
	}

	var peminjaman []models.Peminjaman

	if err := database.DBConn.
		Preload("User").
		Preload("DetailPeminjaman").
		Preload("Room").
		Preload("Pengembalian").
		Where("id_user = ?", userID).
		Find(&peminjaman).
		Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error retrieving data")
	}

	for i := range peminjaman {
		peminjaman[i].TglAcaraFormatted = peminjaman[i].TglAcara.Format("02-01-2006 jam 15:04")
		peminjaman[i].TglAkhirAcaraFormatted = peminjaman[i].TglAkhirAcara.Format("02-01-2006 jam 15:04")
		peminjaman[i].TglAkhirAcaraFormatted = peminjaman[i].TglAkhirAcara.Format("02-01-2006 jam 15:04")

		peminjaman[i].TglAcaraDay = helper.GetIndonesianDay(peminjaman[i].TglAcara.Weekday())
		peminjaman[i].TglAkhirAcaraDay = helper.GetIndonesianDay(peminjaman[i].TglAkhirAcara.Weekday())

		if peminjaman[i].Pengembalian != nil {
			peminjaman[i].Pengembalian.TglPengembalianFormatted = peminjaman[i].Pengembalian.TglPengembalian.Format("02-01-2006 jam 15:04")
			peminjaman[i].Pengembalian.TglPengembalianDay = helper.GetIndonesianDay(peminjaman[i].Pengembalian.TglPengembalian.Weekday())

			fmt.Printf("Peminjaman ID: %d, Tgl Pengembalian: %s, Hari: %s\n",
				peminjaman[i].IdPeminjaman,
				peminjaman[i].Pengembalian.TglPengembalianFormatted,
				peminjaman[i].Pengembalian.TglPengembalianDay)
		} else {
			fmt.Printf("Peminjaman ID: %d tidak memiliki pengembalian.\n", peminjaman[i].IdPeminjaman)
		}

	}

	floors := c.Locals("floors").([]models.Lantai)
	menus := c.Locals("menus").([]models.Menu)

	// Render the detail page with the room data
	return c.Render("history", fiber.Map{
		"Title":         "History",
		"flash_error":   flashError,
		"flash_success": flashSuccess,
		"Peminjaman":    peminjaman,
		"Floors":        floors,
		"menus":         menus,
		"Name":          userName,
	})
}
