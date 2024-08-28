package main

import (
	"log"
	"rangpol/database"
	"rangpol/middleware"
	"rangpol/router"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
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

	// engine := html.New("./views", ".html")
	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Static("/css", "./css")
	app.Static("/vendor", "./vendor")
	app.Static("/img", "./img")
	app.Static("/vendor", "./vendor")
	app.Static("/js", "./js")
	app.Static("/scss", "./scss")

	// store := session.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	app.Use(logger.New())
	// app.Use(store.Handler())

	router.SetupRouters(app)

	app.Listen(":8082")
}
