package controller

import (
	"log"
	"rangpol/database"
	"rangpol/middleware"
	"rangpol/models"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func PeminjamanFormController(c *fiber.Ctx) error {

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

	var user models.User
	if userID != "" {
		if err := database.DBConn.Where("id_user = ?", userID).Find(&user).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error retrieving users")
		}
	}

	// Ambil pesan flash error jika ada
	flashError := sess.Get("flash_error")
	sess.Delete("flash_error")
	if err := sess.Save(); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error saving session")
	}

	// Get the room ID from the query parameters
	idRoomStr := c.Query("id")
	if idRoomStr == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Room ID is required")
	}

	idRoom, err := strconv.ParseUint(idRoomStr, 10, 32) // Konversi string ke uint
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid Room ID format")
	}

	var rooms []models.Room
	// Fetch the room details from the database along with its associated details
	if err := database.DBConn.Find(&rooms).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error retrieving rooms")
	}

	floors := c.Locals("floors").([]models.Lantai)
	menus := c.Locals("menus").(map[string][]models.Menu)
	// log.Printf("Rooms Data: %+v\n", rooms)

	// Render the detail page with the room data
	return c.Render("peminjaman", fiber.Map{
		"Title":       "Peminjaman Ruangan",
		"Rooms":       rooms,
		"flash_error": flashError,
		"User":        user,
		"Floors":      floors,
		"menus":       menus,
		"Dashboard":   "Peminjaman Ruangan",
		"SelectedID":  uint(idRoom),
		"Name":        userName,
	})
}

func PeminjamanController(c *fiber.Ctx) error {

	sess, err := middleware.GetSessionStore().Get(c)
	if err != nil {
		log.Println("Error getting session:", err)
		return err
	}
	// Ambil data dari form
	idUser, _ := strconv.Atoi(c.FormValue("id_user"))
	idRoom, _ := strconv.Atoi(c.FormValue("id_room"))
	namaKegiatan := c.FormValue("nama_kegiatan")
	tglAcaraStr := c.FormValue("tgl_acara")
	tglAkhirAcaraStr := c.FormValue("tgl_akhir_acara")
	tglAcaraWithTZ := tglAcaraStr + "+07:00"
	tglAkhirAcaraWithTZ := tglAkhirAcaraStr + "+07:00"

	// Parse string to time.Time
	tglAcara, err := time.Parse(time.RFC3339, tglAcaraWithTZ)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid datetime format for tgl_acara: " + err.Error())
	}

	tglAkhirAcara, err := time.Parse(time.RFC3339, tglAkhirAcaraWithTZ)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid datetime format for tgl_akhir_acara: " + err.Error())
	}

	PJ := c.FormValue("PJ")
	PA := c.FormValue("PA")
	PK := c.FormValue("PK")
	nTamu, _ := strconv.Atoi(c.FormValue("n_tamu"))
	sifatAcara := c.FormValue("sifat_acara")
	jenisAcara := c.FormValue("jenis_acara")
	keterangan := c.FormValue("keterangan")

	// Inisialisasi database
	db := database.DBConn

	// Insert into peminjaman
	peminjaman := models.Peminjaman{
		IdUser:        uint(idUser),
		IdRoom:        uint(idRoom),
		NamaKegiatan:  namaKegiatan,
		TglAcara:      tglAcara,
		TglAkhirAcara: tglAkhirAcara,
		Status:        0,
		TglAcc:        time.Time{},
	}

	if err := db.Create(&peminjaman).Error; err != nil {
		log.Println("Error inserting peminjaman:", err)

		sess.Set("flash_error", "TglAkhirAcara cannot be before TglAcara")
		if err := sess.Save(); err != nil {
			log.Println("Error saving session:", err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error saving session")
		}

		referer := c.Get("Referer")
		if referer == "" {
			referer = "/" // default fallback jika Referer tidak tersedia
		}
		return c.Redirect(referer)
	}

	// if err := db.Create(&peminjaman).Error; err != nil {
	// 	return c.Status(fiber.StatusInternalServerError).SendString("Error inserting peminjaman 2")
	// }

	// Ambil ID peminjaman yang baru
	idPeminjaman := peminjaman.IdPeminjaman

	// Insert into detail_acara
	detailAcara := models.DetailPeminjaman{
		IdPeminjaman: uint(idPeminjaman),
		PJ:           PJ,
		PA:           PA,
		PK:           PK,
		NTamu:        uint(nTamu),
		SifatAcara:   sifatAcara,
		JenisAcara:   jenisAcara,
		Keterangan:   keterangan,
	}
	if err := db.Create(&detailAcara).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error inserting detail acara")
	}

	// Insert into pengembalian
	pengembalian := models.Pengembalian{
		IdPeminjaman:    uint(idPeminjaman),
		StatusKembali:   false,
		TglPengembalian: time.Time{},
	}
	if err := db.Create(&pengembalian).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error inserting pengembalian")
	}

	// Set pesan flash sukses

	sess.Set("flash_success", "Peminjaman Ruangan Berhasil Di ajukan")

	// Simpan session
	if err := sess.Save(); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error saving session")
	}

	return c.Redirect("/")
}

func TestController(c *fiber.Ctx) error {

	sess, err := middleware.GetSessionStore().Get(c)
	if err != nil {
		return err
	}

	// Set pesan flash sukses
	sess.Set("flash_success", "Peminjaman Ruangan Berhasil Di ajukan")

	// Simpan session
	if err := sess.Save(); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error saving session")
	}

	return c.Redirect("/")

}
