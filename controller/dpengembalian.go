package controller

import (
	"fmt"
	"log"
	"rangpol/database"
	"rangpol/helper"
	"rangpol/middleware"
	"rangpol/models"
	"time"

	"github.com/gofiber/fiber/v2"
)

// AdminPage handles the /admin route
func DataPengembalianController(c *fiber.Ctx) error {
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
	floors := c.Locals("floors").([]models.Lantai)

	var peminjaman []models.Peminjaman

	if err := database.DBConn.
		Preload("User").
		Preload("Pengembalian").
		Where("dlt = ?", 0).
		Find(&peminjaman).Error; err != nil {
		log.Printf("Error retrieving data: %v", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Error retrieving data")
	}

	for i := range peminjaman {
		peminjaman[i].TglAcaraFormatted = peminjaman[i].TglAcara.Format("02-01-2006 jam 15:04")
		peminjaman[i].TglAkhirAcaraFormatted = peminjaman[i].TglAkhirAcara.Format("02-01-2006 jam 15:04")

		peminjaman[i].TglAcaraDay = helper.GetIndonesianDay(peminjaman[i].TglAcara.Weekday())
		peminjaman[i].TglAkhirAcaraDay = helper.GetIndonesianDay(peminjaman[i].TglAkhirAcara.Weekday())

		if peminjaman[i].Pengembalian.TglPengembalian.IsZero() {
			peminjaman[i].Pengembalian.TglPengembalianFormatted = ""
			peminjaman[i].Pengembalian.TglPengembalianDay = ""
			// Jika TglPengembalian masih kosong
			today := time.Now()
			// Jika pengembalian masih kosong, cek tgl akhir acara
			if peminjaman[i].TglAkhirAcara.Before(today) {
				// Hitung selisih hari jika terlambat
				selisihHari := today.Sub(peminjaman[i].TglAkhirAcara).Hours() / 24
				peminjaman[i].KeteranganTelat = fmt.Sprintf("Telat %d hari", int(selisihHari))
			} else {
				// Hitung berapa hari lagi sampai pengembalian
				selisihHari := peminjaman[i].TglAkhirAcara.Sub(today).Hours() / 24
				peminjaman[i].KeteranganTelat = fmt.Sprintf("%d hari lagi", int(selisihHari))
			}
		} else {
			peminjaman[i].Pengembalian.TglPengembalianFormatted = peminjaman[i].Pengembalian.TglPengembalian.Format("02-01-2006 jam 15:04")
			peminjaman[i].Pengembalian.TglPengembalianDay = helper.GetIndonesianDay(peminjaman[i].Pengembalian.TglPengembalian.Weekday())
			// Jika pengembalian sudah ada, cek tgl akhir acara
			if peminjaman[i].Pengembalian.TglPengembalian.Before(peminjaman[i].TglAkhirAcara) {
				peminjaman[i].KeteranganTelat = "Selesai"
			} else if peminjaman[i].Pengembalian.TglPengembalian.Equal(peminjaman[i].TglAkhirAcara) {
				peminjaman[i].KeteranganTelat = "Selesai"
			} else {
				// Jika tgl akhir lebih dari tgl pengembalian
				selisihHari := peminjaman[i].Pengembalian.TglPengembalian.Sub(peminjaman[i].TglAkhirAcara).Hours() / 24

				if int(selisihHari) == 1 {
					// Jika terlambat 1 hari, status tetap oke
					peminjaman[i].KeteranganTelat = "Selesai"
				} else if int(selisihHari) > 1 {
					// Jika terlambat lebih dari 1 hari
					peminjaman[i].KeteranganTelat = fmt.Sprintf("Telat %d hari", int(selisihHari)-1) // -1 karena 1 hari masih oke
				} else {
					peminjaman[i].KeteranganTelat = "Selesai"
				}
			}
		}
	}

	// Set header untuk menonaktifkan caching pada halaman ini
	c.Set("Cache-Control", "no-store")
	c.Set("Pragma", "no-cache")
	c.Set("Expires", "0")

	// Render halaman dengan data yang diperlukan
	fmt.Println("menus : ", floors)

	// add := helper.add()

	return c.Render("datapengembalian", fiber.Map{
		"isIndex":       1,
		"Dashboard":     "Data Pengembalian",
		"flash_error":   flashError,
		"flash_success": flashSuccess,
		"Title":         "Data Pengembalian",
		"menus":         menus,
		"Name":          userName,
		"Peminjaman":    peminjaman,
		"Floors":        floors,
	})
}
