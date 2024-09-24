package middleware

import (
	"fmt"
	"log"
	"rangpol/database"
	"rangpol/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetLantai(c *fiber.Ctx) error {

	var floors []models.Lantai
	if err := database.DBConn.Order("no_lantai ASC").Find(&floors).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error retrieving floors")
	}

	// Set data floors ke context
	c.Locals("floors", floors)

	return c.Next()
}

func GetMenu(c *fiber.Ctx) error {
	sess, err := GetSessionStore().Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error retrieving session")
	}

	roleID := sess.Get("role_id")

	var menus []models.Menu
	if err := database.DBConn.
		Preload("Children", func(db *gorm.DB) *gorm.DB {
			return db.Joins("JOIN usermenus ON usermenus.id_menu = menus.id_menu").
				Where("usermenus.id_akses = ?", roleID).
				Order("menus.urutan ASC")
		}).
		Joins("JOIN usermenus ON usermenus.id_menu = menus.id_menu").
		Where("usermenus.id_akses = ? and ishide=0", roleID).
		Order("menus.urutan ASC").
		Find(&menus).Error; err != nil {
		log.Println("Error retrieving menus:", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Error retrieving menus")
	}

	c.Locals("menus", menus)
	return c.Next()
}

func CheckPrivileges(action string, menuID string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Retrieve user session (example)
		sess, err := GetSessionStore().Get(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error retrieving session")
		}

		userID := sess.Get("user_id")
		roleID := sess.Get("role_id")

		if userID == nil || roleID == nil {
			return c.Redirect("/login")
		}

		// Check if the user has the right privileges
		hasAccess := getPrivileges(database.DBConn, action, menuID, roleID.(uint))

		userName := sess.Get("name_user")
		menus := c.Locals("menus").([]models.Menu)
		floors := c.Locals("floors").([]models.Lantai)

		if !hasAccess {
			// return c.Status(fiber.StatusForbidden).SendString("You do not have access to this menu")
			return c.Render("403", fiber.Map{
				"menus":  menus,
				"Floors": floors,
				"Name":   userName,
			})
		}

		// Continue if user has access
		return c.Next()
	}
}

func getPrivileges(db *gorm.DB, param string, menuID string, roleID uint) bool {
	var usermenu models.Usermenu

	// Query user privileges for the specific menu
	err := db.Where("id_akses = ? AND id_menu = ?", roleID, menuID).First(&usermenu).Error
	if err != nil {
		fmt.Println("Error fetching privileges:", err)
		return false
	}

	// Check the specific action (view, create, edit, del, print)
	switch param {
	case "view":
		return usermenu.View == 1
	case "create":
		return usermenu.Add == 1
	case "edit":
		return usermenu.Edit == 1
	case "del":
		return usermenu.Delete == 1
	case "print":
		return usermenu.Print == 1
	default:
		return false
	}
}
