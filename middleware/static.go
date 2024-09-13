package middleware

import (
	"fmt"
	"log"
	"rangpol/database"
	"rangpol/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func LantaiMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Ambil data menu dari database
		var floors []models.Lantai
		if err := database.DBConn.Order("no_lantai ASC").Find(&floors).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error retrieving floors")
		}

		// Set data menu ke context
		c.Locals("floors", floors)

		return c.Next()
	}
}

func GetMenu(c *fiber.Ctx) error {
	var menus []models.Menu
	if err := database.DBConn.Preload("Children").Order("urutan ASC").Find(&menus).Order("Urutan asc").Error; err != nil {
		log.Println("Error retrieving menus:", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Error retrieving menus")
	}

	// log.Println("Grouped Menus:", groupedMenus)
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

		if !hasAccess {
			return c.Status(fiber.StatusForbidden).SendString("You do not have access to this menu")
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
