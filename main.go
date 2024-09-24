package main

import (
	"log"
	"rangpol/database"
	"rangpol/middleware"
	"rangpol/router"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
)

func init() {

	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Eror in loading .env file")
	}

	database.ConnectDB()
	middleware.InitSessionStore()
}

func main() {
	sqlDB, err := database.DBConn.DB()
	if err != nil {
		log.Fatalf("Error in SQL connection: %v", err)
	}
	defer sqlDB.Close()

	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views:             engine,
		BodyLimit:         5 * 1024 * 1024, // Batas body 5 MB
		StreamRequestBody: true,
	})

	// Middleware
	// Middleware untuk menangani recover dari panic
	app.Use(recover.New())

	// Error handler
	app.Use(func(c *fiber.Ctx) error {
		if err := c.Next(); err != nil {
			if fiberErr, ok := err.(*fiber.Error); ok {
				// Cek apakah error adalah akibat body limit
				if fiberErr.Code == fiber.StatusRequestEntityTooLarge {
					return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
						"error": "File terlalu besar. Batas ukuran maksimum adalah 10 MB.",
					})
				}
			}
			return err // Kembalikan error lain
		}
		return nil
	}) // Menambahkan middleware recover
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))
	app.Use(logger.New())
	app.Use(middleware.GetLantai)
	app.Use(middleware.GetMenu)
	app.Use(middleware.RedirectIfAuthenticated)

	// Static files
	app.Static("/css", "./views/css")
	app.Static("/vendor", "./views/vendor")
	app.Static("/img", "./views/img")
	app.Static("/js", "./views/js")

	router.SetupRouters(app)

	if err := app.Listen(":8082"); err != nil {
		log.Fatal(err)
	}
}
