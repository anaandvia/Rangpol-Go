package database

import (
	"log"
	"os"
	"rangpol/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DBConn *gorm.DB

func ConnectDB() {

	user := os.Getenv("db_user")
	password := os.Getenv("db_password")
	dbname := os.Getenv("db_name")

	dsn := user + ":" + password + "@tcp(127.0.0.1:3306)/" + dbname + "?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})

	if err != nil {
		panic("Database connection failed.")
	}

	log.Println("Connection successful.")

	// auto migration
	// Auto migration for multiple models
	db.AutoMigrate(&models.Menu{}, &models.Usermenu{}, &models.User{}, &models.Blog{}, &models.Room{}, &models.Lantai{}, &models.DetailRoom{}, &models.Peminjaman{}, &models.DetailPeminjaman{}, &models.Fasilitas{}, &models.FasilitasRoom{}, &models.Jurusan{}, &models.Pengembalian{})

	DBConn = db
}
