package router

import (
	"rangpol/controller"
	"rangpol/middleware"

	"github.com/gofiber/fiber/v2"
)

// setup routing information
func SetupRouters(app *fiber.App) {

	//list => get
	//add => post
	//update => put
	//delete => delete

	// ------------------------index-------------------------------
	app.Get("/", controller.HomeController)
	// ------------------------------------------------------------
	// --------------------- Registration -------------------------
	app.Get("/register", controller.RegisterFormController)
	app.Post("/register", controller.RegisterController)
	// ------------------------------------------------------------
	// ------------------ authentication ---------------------------
	app.Get("/login", middleware.RedirectIfAuthenticated, controller.LoginFormController)
	app.Post("/login", controller.LoginPostController)
	app.Get("/logout", controller.LogoutController)
	// ------------------------------------------------------------
	// ------------------ Peminjaman ---------------------------
	app.Get("/peminjaman", middleware.CheckPrivileges("view", "7"), controller.PeminjamanFormController)
	app.Post("/borang", middleware.CheckPrivileges("add", "7"), controller.PeminjamanController)
	// ------------------------------------------------------------
	// ------------------ Ruangan ---------------------------
	app.Get("/detail_room", middleware.CheckPrivileges("view", "2"), controller.RoomDetailController)
	// ------------------------------------------------------------
	// ------------------ Pengembalian ---------------------------
	app.Get("/kembalikan", middleware.CheckPrivileges("view", "8"), controller.PengembalianFormController)
	app.Post("/kembalikan", middleware.CheckPrivileges("edit", "8"), controller.PengembalianController)
	// ------------------------------------------------------------

	// app.Get("/peminjaman", controller.PeminjamanFormController)

}
