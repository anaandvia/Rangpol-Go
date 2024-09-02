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
	// log.Printf("Rooms Data: %+v\n", rooms)

	// Render the detail page with the room data
	return c.Render("peminjaman", fiber.Map{
		"Title":       "Peminjaman Ruangan",
		"Rooms":       rooms,
		"flash_error": flashError,
		"User":        user,
		"Floors":      floors,
		"Dashboard":   "Peminjaman Ruangan",
		"SelectedID":  uint(idRoom),
	})
}

const datetimeLayout = "2006-01-02T15:04"

func PeminjamanController(c *fiber.Ctx) error {
	// Ambil data dari form
	idUser, _ := strconv.Atoi(c.FormValue("id_user"))
	idRoom, _ := strconv.Atoi(c.FormValue("id_room"))
	namaKegiatan := c.FormValue("nama_kegiatan")
	tglAcaraStr := c.FormValue("tgl_acara")
	tglAkhirAcaraStr := c.FormValue("tgl_akhir_acara")
	// Parse string to time.Time
	tglAcara, err := time.Parse(datetimeLayout, tglAcaraStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid datetime format for tgl_acara: " + err.Error())
	}

	tglAkhirAcara, err := time.Parse(datetimeLayout, tglAkhirAcaraStr)
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
	log.Println(uint(idRoom))
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
		return c.Status(fiber.StatusInternalServerError).SendString("Error inserting peminjaman")
	}

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

	log.Println("Data successfully inserted")

	// Redirect to index page
	// Ambil pesan flash error jika ada

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