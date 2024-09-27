package controller

import (
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"rangpol/database"
	"rangpol/helper"
	"rangpol/middleware"
	"rangpol/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func DataRuanganController(c *fiber.Ctx) error {
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

	var rooms []models.Room
	if err := database.DBConn.Where("dlt = ?", 0).Find(&rooms).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error retrieving floors")
	}

	// Set header untuk menonaktifkan caching pada halaman ini
	c.Set("Cache-Control", "no-store")
	c.Set("Pragma", "no-cache")
	c.Set("Expires", "0")

	// Render halaman dengan data yang diperlukan
	fmt.Println("menus : ", floors)

	// add := helper.add()

	return c.Render("dataruangan", fiber.Map{
		"isIndex":       1,
		"Dashboard":     "Data Ruangan",
		"flash_error":   flashError,
		"flash_success": flashSuccess,
		"Title":         "Data Ruangan",
		"menus":         menus,
		"Name":          userName,
		"Room":          rooms,
		"Floors":        floors,
	})
}

func CreateRuanganController(c *fiber.Ctx) error {

	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in create:", r)
			c.Status(fiber.StatusInternalServerError).SendString("Internal server error")
		}
	}()

	sess, err := middleware.GetSessionStore().Get(c)
	if err != nil {
		log.Println("Error getting session:", err)
		return err
	}

	// Retrieve form values
	noRoom := c.FormValue("no_room")
	nameRoom := c.FormValue("name_room")
	lantai, _ := strconv.Atoi(c.FormValue("lantai"))
	kapasitas, _ := strconv.Atoi(c.FormValue("kapasitas"))
	status, _ := strconv.ParseBool(c.FormValue("status"))
	file, err := c.FormFile("foto")

	var fileNameOnly string

	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Error getting file: " + err.Error())
	}

	if file != nil {
		const maxFileSize = 10 * 1024 * 1024 // 10 MB
		if file.Size > maxFileSize {
			sess.Set("flash_error", "File too large. Max size allowed is 10 MB.")
			if err := sess.Save(); err != nil {
				return c.Status(fiber.StatusInternalServerError).SendString("Error saving session")
			}
			return c.Redirect("/admin/dataruangan")
		}

		if !helper.IsValidFileType(file) {
			sess.Set("flash_error", "Invalid file type. Only JPG, PNG, and JPEG are allowed.")
			if err := sess.Save(); err != nil {
				return c.Status(fiber.StatusInternalServerError).SendString("Error saving session")
			}
			return c.Redirect("/admin/dataruangan")
		}

		// Define the file upload path
		uploadPath := "./views/img/ruangan/" // Make sure this directory exists and is writable
		newFileName := helper.RenameFile(file.Filename)

		// Gabungkan path dengan nama file baru
		filePath := filepath.Join(uploadPath, newFileName)
		fileNameOnly = newFileName

		log.Println(fileNameOnly)

		// Save the file
		if err := c.SaveFile(file, filePath); err != nil {
			sess.Set("flash_error", "Error saving file: "+err.Error())
			if err := sess.Save(); err != nil {
				return c.Status(fiber.StatusInternalServerError).SendString("Error saving session")
			}
			return c.Redirect("/admin/dataruangan")
		}
	}

	// Check if room number already exists
	var existingRoom models.Room
	if err := database.DBConn.Where("no_room = ?", noRoom).First(&existingRoom).Error; err == nil {
		sess.Set("flash_error", "Room number already in use")
		if err := sess.Save(); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error saving session")
		}
		return c.Redirect("/admin/dataruangan")
	}

	// Create a new room entry
	room := models.Room{
		No_room:   noRoom,
		Name_room: nameRoom,
		Lantai:    uint(lantai),
		Kapasitas: uint(kapasitas),
		Foto:      fileNameOnly, // Set the saved file path
		Status:    status,
	}

	log.Println("apaaaa  " + fileNameOnly)

	// Save the new room to the database
	if err := database.DBConn.Create(&room).Error; err != nil {
		sess.Set("flash_error", "Error creating room")
		if err := sess.Save(); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error saving session")
		}
		return c.Status(fiber.StatusInternalServerError).SendString("Error creating room")
	}

	// Set success flash message
	sess.Set("flash_success", "Room created successfully")

	// Save session
	if err := sess.Save(); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error saving session")
	}

	return c.Redirect("/admin/dataruangan")
}

func UpdateRuanganController(c *fiber.Ctx) error {

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
	roomID := c.Params("id")

	// Retrieve form values
	noRoom := c.FormValue("no_room")
	nameRoom := c.FormValue("name_room")
	lantai, _ := strconv.Atoi(c.FormValue("lantai"))
	kapasitas, _ := strconv.Atoi(c.FormValue("kapasitas"))
	status, _ := strconv.ParseBool(c.FormValue("status"))
	oldFoto := c.FormValue("oldfoto")
	log.Println("test")
	file, err := c.FormFile("foto")

	// File handling (for foto)
	// file, err := c.FormFile("foto")
	// s := file
	// log.Println(s)

	var fileNameOnly string

	if err != nil {
		if file != nil {
			// Jika tidak ada file yang diunggah, mungkin ini bukan error kritis
			// Tangani dengan cara lain, seperti menggunakan file lama
			return c.Status(fiber.StatusBadRequest).SendString("Error getting file: " + err.Error())
		}

		fileNameOnly = oldFoto
	}

	if file != nil {
		const maxFileSize = 5 * 1024 * 1024 // 10 MB
		if file.Size > maxFileSize {
			// return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("File too large: %s. Max size allowed is 10 MB.", file.Filename))
			sess.Set("flash_error", "File too large. Max size allowed is 10 MB.")
			if err := sess.Save(); err != nil {
				return c.Status(fiber.StatusInternalServerError).SendString("Error saving session")
			}
			return c.Redirect("/admin/dataruangan")
		}

		if !helper.IsValidFileType(file) {
			// return c.Status(fiber.StatusInternalServerError).SendString("Invalid file type. Only JPG, PNG, and JPEG are allowed.")
			sess.Set("flash_error", "Invalid file type. Only JPG, PNG, and JPEG are allowed.")
			if err := sess.Save(); err != nil {
				return c.Status(fiber.StatusInternalServerError).SendString("Error saving session")
			}
			return c.Redirect("/admin/dataruangan")
		}

		// Define the file upload path
		uploadPath := "./views/img/ruangan/" // Make sure this directory exists and is writable

		newFileName := helper.RenameFile(file.Filename)

		// Gabungkan path dengan nama file baru
		filePath := filepath.Join(uploadPath, newFileName)

		fileNameOnly = newFileName

		log.Println(fileNameOnly)

		// Save the file
		if err := c.SaveFile(file, filePath); err != nil {
			sess.Set("flash_error", "Error saving file: "+err.Error())
			if err := sess.Save(); err != nil {
				return c.Status(fiber.StatusInternalServerError).SendString("Error saving session")
			}
			return c.Redirect("/admin/dataruangan")
			// return c.Status(fiber.StatusInternalServerError).SendString("Error saving file: " + err.Error())
		}
	} else {
		fileNameOnly = oldFoto
	}

	// Check if room number is already used by another room
	var existingRoom models.Room
	if err := database.DBConn.Where("no_room = ? AND id_room != ?", noRoom, roomID).First(&existingRoom).Error; err == nil {
		// Room number is already in use
		sess.Set("flash_error", "Room number already in use")
		if err := sess.Save(); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error saving session")
		}
		return c.Redirect("/admin/dataruangan")
	}

	// Find the room by ID
	var room models.Room
	if err := database.DBConn.First(&room, roomID).Error; err != nil {
		sess.Set("flash_error", "Room not found")
		if err := sess.Save(); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error saving session")
		}
		return c.Redirect("/admin/dataruangan")
	}

	// Update room details
	room.No_room = noRoom
	room.Name_room = nameRoom
	room.Lantai = uint(lantai)
	room.Kapasitas = uint(kapasitas)
	room.Foto = fileNameOnly // Set the saved file path
	room.Status = status

	log.Println(room.Foto)

	// Save the updated room to the database
	if err := database.DBConn.Save(&room).Error; err != nil {
		sess.Set("flash_error", "Error updating room")
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

	return c.Redirect("/admin/dataruangan")
}

func DeleteRuanganController(c *fiber.Ctx) error {

	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in delete:", r)
			c.Status(fiber.StatusInternalServerError).SendString("Internal server error")
		}
	}()

	sess, err := middleware.GetSessionStore().Get(c)
	if err != nil {
		log.Println("Error getting session:", err)
		return err
	}

	// Retrieve room ID from the URL path
	roomID := c.Params("id")

	// Find the room by ID
	var room models.Room
	if err := database.DBConn.First(&room, roomID).Error; err != nil {
		sess.Set("flash_error", "Room not found")
		if err := sess.Save(); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error saving session")
		}
		return c.Redirect("/admin/dataruangan")
	}

	var existingPeminjaman models.Peminjaman
	if err := database.DBConn.Where("id_room = ? and dlt = ?", roomID, 0).First(&existingPeminjaman).Error; err == nil {
		// Room number is already in use
		sess.Set("flash_error", "Tidak dapat di hapus! Ruangan digunakan pada peminjaman ruangan")
		if err := sess.Save(); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error saving session")
		}
		return c.Redirect("/admin/dataruangan")
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		// Jika terjadi error selain data tidak ditemukan, maka kembalikan error lainnya
		fmt.Println("Error: ", err)

		// Tampilkan pesan error umum kepada pengguna
		return c.Status(fiber.StatusInternalServerError).SendString("Terjadi kesalahan pada database, silakan coba lagi.")
	}

	// Update the dlt status to indicate the room is deleted
	room.Dlt = 1 // Assuming Dlt is the field that indicates deletion

	// Save the updated room to the database
	if err := database.DBConn.Save(&room).Error; err != nil {
		sess.Set("flash_error", "Error marking room as deleted")
		if err := sess.Save(); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error saving session")
		}
		return c.Status(fiber.StatusInternalServerError).SendString("Error marking room as deleted")
	}

	// Set success flash message
	sess.Set("flash_success", "Room deleted successfully")

	// Save session
	if err := sess.Save(); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error saving session")
	}

	return c.Redirect("/admin/dataruangan")
}
