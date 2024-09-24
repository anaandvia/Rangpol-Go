package helper

import (
	"fmt"
	"log"
	"math/rand"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"
)

func GetIndonesianDay(day time.Weekday) string {
	days := map[time.Weekday]string{
		time.Sunday:    "Minggu",
		time.Monday:    "Senin",
		time.Tuesday:   "Selasa",
		time.Wednesday: "Rabu",
		time.Thursday:  "Kamis",
		time.Friday:    "Jumat",
		time.Saturday:  "Sabtu",
	}
	return days[day]
}

func toUpperCase(s interface{}) string {
	// Assert the type to string
	if str, ok := s.(string); ok {
		return strings.ToUpper(str)
	}
	return ""
}

func Add(a, b int) int {
	return a + b
}

// Validasi MIME type file
func IsValidFileType(file *multipart.FileHeader) bool {
	// Mendapatkan MIME type file
	fileType := file.Header.Get("Content-Type")

	log.Println(fileType)

	// Daftar MIME type yang diizinkan
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/jpg":  true,
	}

	return allowedTypes[fileType]
}

func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// Fungsi untuk mengubah nama file
func RenameFile(originalName string) string {
	// Dapatkan ekstensi file dan ubah menjadi huruf kecil
	extension := strings.ToLower(filepath.Ext(originalName))

	// Dapatkan nama file tanpa ekstensi
	baseName := originalName[:len(originalName)-len(extension)]

	// Buat string acak yang berisi angka dan huruf
	randomStr := RandomString(8) // Misal panjangnya 8 karakter

	// Gabungkan nama file baru dengan string acak dan ekstensi file (yang sudah huruf kecil)
	newName := fmt.Sprintf("%s_%s%s", baseName, randomStr, extension)

	return newName
}
