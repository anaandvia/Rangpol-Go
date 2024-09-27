package controller

import (
	"fmt"
	"log"
	"path/filepath"
	"rangpol/database"
	"rangpol/helper"
	_ "rangpol/helper"
	"rangpol/middleware"
	"rangpol/models"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func formatDateTime(t time.Time) string {
	return t.Format("2006-01-02T15:04:05")
}

func PengembalianFormController(c *fiber.Ctx) error {

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
	log.Println("log: ", flashError)
	// Get the room ID from the query parameters
	idPeminjamanStr := c.Query("id")
	if idPeminjamanStr == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Peminjaman ID is required")
	}

	idPeminjaman, err := strconv.ParseUint(idPeminjamanStr, 10, 32) // Konversi string ke uint
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid Peminjaman ID format")
	}

	var peminjaman models.Peminjaman

	if err := database.DBConn.
		Preload("User").
		Preload("DetailPeminjaman").
		Preload("Room").
		Preload("Pengembalian").
		Where("id_peminjaman = ?", idPeminjaman).
		Find(&peminjaman).
		Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error retrieving data")
	}

	formattedTglAcara := formatDateTime(peminjaman.TglAcara)
	formattedTglPengembalian := formatDateTime(peminjaman.Pengembalian.TglPengembalian)
	statuskembali := 0
	if peminjaman.Pengembalian.StatusKembali == 1 {
		statuskembali = 1
	}

	floors := c.Locals("floors").([]models.Lantai)
	menus := c.Locals("menus").([]models.Menu)
	// log.Printf("Peminjaman Data: %+v\n", peminjaman)

	// Render the detail page with the Peminjaman data
	return c.Render("pengembalian", fiber.Map{
		"Title":           "Pengembalian Ruangan",
		"Peminjaman":      peminjaman,
		"flash_error":     flashError,
		"Floors":          floors,
		"menus":           menus,
		"Dashboard":       "Pengembalian Ruangan",
		"Name":            userName,
		"TglAcara":        formattedTglAcara,
		"TglPengembalian": formattedTglPengembalian,
		"SKembali":        statuskembali,
	})
}

func PengembalianController(c *fiber.Ctx) error {
	log.Println("dipanggil")
	sess, err := middleware.GetSessionStore().Get(c)
	if err != nil {
		log.Println("Error getting session:", err)
		return err
	}
	// Ambil data dari form
	idPengembalian, _ := strconv.Atoi(c.FormValue("id_pengembalian"))
	idPeminjaman, _ := strconv.Atoi(c.FormValue("id_peminjaman"))
	kendala := c.FormValue("kendala")
	tglPengembalianStr := c.FormValue("tgl_pengembalian")
	tglPengembalianWithTZ := tglPengembalianStr + "+07:00"

	// Parse string to time.Time
	tglPengembalian, err := time.Parse(time.RFC3339, tglPengembalianWithTZ)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid datetime format for tgl_pengembalian: " + err.Error())
	}

	file, err := c.FormFile("foto_b")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Error getting file: " + err.Error())
	}

	if !helper.IsValidFileType(file) {
		// return c.Status(fiber.StatusInternalServerError).SendString("Invalid file type. Only JPG, PNG, and JPEG are allowed.")
		sess.Set("flash_error", "Invalid file type. Only JPG, PNG, and JPEG are allowed.")
		if err := sess.Save(); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error saving session")
		}
		return c.Redirect(fmt.Sprintf("/kembalikan?id=%d", idPeminjaman))

		// referer := c.Get("Referer")
		// if referer == "" {
		// 	referer = "/history" // default fallback jika Referer tidak tersedia
		// }
		// return c.Redirect(referer)
	}

	// Define the file upload path
	uploadPath := "./views/img/pengembalian/" // Make sure this directory exists and is writable

	newFileName := helper.RenameFile(file.Filename)

	// Gabungkan path dengan nama file baru
	filePath := filepath.Join(uploadPath, newFileName)

	// Save the file
	if err := c.SaveFile(file, filePath); err != nil {
		sess.Set("flash_error", "Error saving file: "+err.Error())
		if err := sess.Save(); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error saving session")
		}
		return c.Redirect("/kembalikan?id=?", idPeminjaman)
		// return c.Status(fiber.StatusInternalServerError).SendString("Error saving file: " + err.Error())
	}

	db := database.DBConn

	pengembalian := models.Pengembalian{
		FotoB:           newFileName,
		Kendala:         kendala,
		TglPengembalian: tglPengembalian,
		StatusKembali:   1,
	}
	// Update Pengembalian record
	if err := db.Where("id = ?", idPengembalian).Updates(&pengembalian).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error inserting pengembalian")
	}

	// Set pesan flash sukses
	sess.Set("flash_success", "Pengembalian Ruangan Berhasil Di ajukan")
	// Simpan session
	if err := sess.Save(); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error saving session")
	}

	return c.Redirect("/history")
}
