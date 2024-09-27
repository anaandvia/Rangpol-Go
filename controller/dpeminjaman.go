package controller

import (
	"fmt"
	"log"
	"rangpol/database"
	"rangpol/helper"
	"rangpol/middleware"
	"rangpol/models"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

// AdminPage handles the /admin route
func DataPeminjamanController(c *fiber.Ctx) error {
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

	peminjamanID := c.Query("id")

	var peminjaman []models.Peminjaman
	query := database.DBConn.
		Preload("User").
		Preload("DetailPeminjaman").
		Preload("Room").
		Preload("Pengembalian").
		Where("dlt = ?", 0)

	// Tambahkan kondisi jika peminjamanID ada
	if peminjamanID != "" {
		query = query.Where("id_peminjaman = ?", peminjamanID) // Ganti "id" dengan nama kolom yang sesuai jika berbeda
	}

	if err := query.Find(&peminjaman).Error; err != nil {
		log.Printf("Error retrieving data: %v", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Error retrieving data")
	}

	for i := range peminjaman {
		peminjaman[i].TglAcaraFormatted = peminjaman[i].TglAcara.Format("02-01-2006 jam 15:04")
		peminjaman[i].TglAkhirAcaraFormatted = peminjaman[i].TglAkhirAcara.Format("02-01-2006 jam 15:04")

		peminjaman[i].TglAcaraDay = helper.GetIndonesianDay(peminjaman[i].TglAcara.Weekday())
		peminjaman[i].TglAkhirAcaraDay = helper.GetIndonesianDay(peminjaman[i].TglAkhirAcara.Weekday())
	}
	// Set header untuk menonaktifkan caching pada halaman ini
	c.Set("Cache-Control", "no-store")
	c.Set("Pragma", "no-cache")
	c.Set("Expires", "0")

	// Render halaman dengan data yang diperlukan
	fmt.Println("menus : ", floors)

	// add := helper.add()

	return c.Render("datapeminjaman", fiber.Map{
		"isIndex":       1,
		"Dashboard":     "Data Peminjaman",
		"flash_error":   flashError,
		"flash_success": flashSuccess,
		"Title":         "Data Peminjaman",
		"menus":         menus,
		"Name":          userName,
		"Peminjaman":    peminjaman,
		"Floors":        floors,
	})
}

func UpdatePeminjamanController(c *fiber.Ctx) error {

	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in upload:", r)
			c.Status(fiber.StatusInternalServerError).SendString("Internal server error")
		}
	}()

	sess, err := middleware.GetSessionStore().Get(c)
	if err != nil {
		log.Println("Error getting session:", err)
		return err
	}

	// Retrieve room ID from the URL path
	peminjamanID := c.Params("id")

	// Retrieve form values
	keterangan := c.FormValue("keterangan")
	status, _ := strconv.Atoi(c.FormValue("status"))

	// Ambil waktu saat ini dalam format RFC3339
	tglAccWithTZ := time.Now().Format(time.RFC3339)

	// Parse waktu dengan format RFC3339 menjadi time.Time
	tglAcc, err := time.Parse(time.RFC3339, tglAccWithTZ)
	if err != nil {
		log.Println("Error parsing time:", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Error parsing time")
	}

	// Find the peminjaman by ID
	var peminjaman models.Peminjaman
	if err := database.DBConn.Where("id_peminjaman = ? and dlt = ?", peminjamanID, 0).First(&peminjaman).Error; err != nil {
		sess.Set("flash_error", "Peminjaman not found")
		if err := sess.Save(); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error saving session")
		}
		return c.Redirect("/admin/datapeminjaman")
	}
	log.Println(tglAcc)
	// Update room details
	peminjaman.Status = uint(status)
	peminjaman.KetVerif = keterangan
	peminjaman.TglAcc = tglAcc

	// Save the updated peminjaman to the database
	if err := database.DBConn.Save(&peminjaman).Error; err != nil {
		log.Println("Error:", err)
		sess.Set("flash_error", "Error updating room: "+err.Error())
		if err := sess.Save(); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error saving session")
		}
		return c.Status(fiber.StatusInternalServerError).SendString("Error updating room")
	}

	// Set success flash message
	sess.Set("flash_success", "Room updated successfully")

	// Save session
	if err := sess.Save(); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error saving session")
	}

	return c.Redirect("/admin/datapeminjaman")
}
