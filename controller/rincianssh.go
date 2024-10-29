package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"rangpol/models"
	_ "rangpol/models"
	"runtime"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func DataRincianSshController(c *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic: %v", r)
			logMemoryUsage()
		}
	}()

	resp, err := http.Get("http://localhost:8080/sipblud/sipblud-all-tukd-dev/index.php/api/rincianssh") // Ganti dengan URL yang sesuai
	if err != nil {
		fmt.Println("Error:", err)
	}
	// fmt.Println(resp)
	defer resp.Body.Close()

	var rinciansshList []models.Rincianssh
	if err := json.NewDecoder(resp.Body).Decode(&rinciansshList); err != nil {
		fmt.Println("Error decoding JSON:", err)
	}
	// for _, r := range rinciansshList {
	// 	fmt.Printf("ID: %s, Name: %s, Price: %s\n", r.RinciansshID, r.RinciansshNama, r.Harga)
	// }
	logMemoryUsage()
	fmt.Printf("Fetched %d records\n", len(rinciansshList))

	return c.Render("rincianssh", fiber.Map{
		"Rincianssh": rinciansshList,
	})
}

func logMemoryUsage() {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	log.Printf("Memory Usage: Alloc = %v MiB, TotalAlloc = %v MiB, Sys = %v MiB",
		memStats.Alloc/1024/1024,
		memStats.TotalAlloc/1024/1024,
		memStats.Sys/1024/1024)
}

func fetchPaginatedRincianssh(db *gorm.DB, page int, pageSize int) ([]models.Rincianssh, error) {
	var rinciansshList []models.Rincianssh
	offset := (page - 1) * pageSize
	if err := db.Offset(offset).Limit(pageSize).Find(&rinciansshList).Error; err != nil {
		return nil, err
	}
	return rinciansshList, nil
}

func GetRincianssh(c *fiber.Ctx) error {
	var rincianssh []models.Rincianssh
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize", "20"))
	searchText := c.Query("searchtext", "")
	filterId := c.Query("filterssrinciansshid", "")

	db := c.Locals("db").(*gorm.DB) // Ambil koneksi DB dari konteks Fiber

	// Mulai query
	query := db.Model(&models.Rincianssh{}).
		Select("rincianssh.*, tr.ssrinciansshkode, tr.ssrinciansshnama").
		Joins("LEFT JOIN tbmssrincianssh tr ON tr.ssrinciansshid = rincianssh.ssrinciansshid AND tr.dlt = ?", false).
		Where("rincianssh.dlt = ? AND rincianssh.thang = ?", false, "2024") // Ubah sesuai tahun yang relevan

	// Tambahkan kondisi pencarian
	if searchText != "" {
		query = query.Where("rincianssh.rinciansshkode ILIKE ? OR rincianssh.rinciansshnama ILIKE ?", "%"+searchText+"%", "%"+searchText+"%")
	}

	// Tambahkan filter ID
	if filterId != "" && filterId != "-" {
		query = query.Where("rincianssh.ssrinciansshid = ?", filterId)
	}

	// Pagination
	offset := (page - 1) * pageSize
	query = query.Offset(offset).Limit(pageSize)

	// Urutan default
	query = query.Order("LPAD(tr.ssrinciansshkode, 3, '0'), LPAD(rincianssh.rinciansshkode, 3, '0')")

	// Eksekusi query
	if err := query.Find(&rincianssh).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Menghitung total data
	var totalRecords int64
	db.Model(&models.Rincianssh{}).Count(&totalRecords)

	return c.JSON(fiber.Map{
		"data":       rincianssh,
		"page":       page,
		"pageSize":   pageSize,
		"totalCount": totalRecords,
	})
}

func GetRinciansshFromAPI(c *fiber.Ctx) error {
	// 1. Ambil data dari API eksternal
	resp, err := http.Get("http://localhost:8080/sipblud/sipblud-all-tukd-dev/index.php/api/rincianssh")
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil data dari API eksternal"})
	}
	defer resp.Body.Close()

	// 2. Decode JSON response
	var allData []models.Rincianssh
	if err := json.NewDecoder(resp.Body).Decode(&allData); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal decode data JSON"})
	}

	// 3. Implementasi Pagination
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize", "20"))

	start := (page - 1) * pageSize
	end := start + pageSize

	if start > len(allData) {
		return c.JSON(fiber.Map{"data": []models.Rincianssh{}, "page": page, "pageSize": pageSize, "totalCount": len(allData)})
	}
	if end > len(allData) {
		end = len(allData)
	}

	// 4. Ambil data berdasarkan halaman
	paginatedData := allData[start:end]

	// 5. Return data dengan informasi pagination
	return c.JSON(fiber.Map{
		"data":       paginatedData,
		"page":       page,
		"pageSize":   pageSize,
		"totalCount": len(allData),
	})
}

func GetPaginatedRincianssh(c *fiber.Ctx, db *gorm.DB) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize", "10"))

	rinciansshList, total, err := models.GetPaginatedRincianssh(db, page, pageSize)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Gagal mengambil data")
	}

	// Mengirim data beserta total untuk frontend
	return c.JSON(fiber.Map{
		"data":       rinciansshList,
		"total":      total,
		"page":       page,
		"pageSize":   pageSize,
		"totalPages": (total + int64(pageSize) - 1) / int64(pageSize),
	})
}

func DataRincianSsh1Controller(c *fiber.Ctx, db *gorm.DB) error {
	// Baca parameter `page` dan `pageSize` dari query string
	page, err := strconv.Atoi(c.Query("page", "1")) // Default ke halaman 1
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.Query("pageSize", "10")) // Default ke 10 item per halaman
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	var rinciansshList []models.Rincianssh
	offset := (page - 1) * pageSize

	// Query dengan limit dan offset untuk pagination
	if err := db.Offset(offset).Limit(pageSize).Find(&rinciansshList).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Gagal mengambil data",
		})
	}

	return c.JSON(fiber.Map{
		"data":      rinciansshList,
		"page":      page,
		"pageSize":  pageSize,
		"totalData": len(rinciansshList),
	})
}
