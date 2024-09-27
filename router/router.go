package router

import (
	"rangpol/controller"
	"rangpol/middleware"

	"github.com/gofiber/fiber/v2"
)

// setup routing information
func SetupRouters(app *fiber.App) {

	//list => get
	//create => post
	//update => put
	//delete => delete
	// ------------------------index-------------------------------
	app.Get("/", middleware.CheckPrivileges("view", "1"), controller.HomeController)
	app.Get("/admin", middleware.CheckPrivileges("view", "9"), controller.AdminPage)

	// ------------------------------------------------------------
	// --------------------- Registration -------------------------
	app.Get("/register", controller.RegisterFormController)
	app.Post("/register", controller.RegisterController)
	// ------------------------------------------------------------
	// ------------------ authentication ---------------------------
	app.Get("/login", controller.LoginFormController)
	app.Post("/login", controller.LoginPostController)
	app.Get("/logout", controller.LogoutController)
	// ------------------------------------------------------------
	// ------------------ Peminjaman ---------------------------
	app.Get("/peminjaman", middleware.CheckPrivileges("view", "7"), controller.PeminjamanFormController)
	app.Post("/borang", middleware.CheckPrivileges("create", "7"), controller.PeminjamanController)
	// ------------------------------------------------------------
	// ------------------ Ruangan ---------------------------------
	app.Get("/detail_room", middleware.CheckPrivileges("view", "2"), controller.RoomDetailController)
	// ------------------------------------------------------------
	// ------------------ Pengembalian ---------------------------
	app.Get("/kembalikan", middleware.CheckPrivileges("view", "8"), controller.PengembalianFormController)
	app.Post("/kembalikan", middleware.CheckPrivileges("edit", "8"), controller.PengembalianController)
	// ------------------------------------------------------------
	// ------------------ History ---------------------------------
	app.Get("/history", middleware.CheckPrivileges("view", "5"), controller.HistoryPeminjamanController)
	// ------------------------------------------------------------

	// ------------------Superadmin---------------------------------
	// -------------------------- Data User ------------------------
	app.Get("/admin/datauser", middleware.CheckPrivileges("view", "11"), controller.DataUserController)
	app.Post("/admin/datauser/tambah", middleware.CheckPrivileges("create", "11"), controller.CreateUserController)
	app.Post("/admin/datauser/delete/:id", middleware.CheckPrivileges("del", "11"), controller.DeleteUserController)
	app.Post("/admin/datauser/update/:id", middleware.CheckPrivileges("edit", "11"), controller.UpdateUserController)
	// -------------------------------------------------------------
	// ----------------------- Data Ruangan ------------------------
	app.Get("/admin/dataruangan", middleware.CheckPrivileges("view", "9"), controller.DataRuanganController)
	app.Post("/admin/dataruangan/tambah", middleware.CheckPrivileges("create", "9"), controller.CreateRuanganController)
	app.Post("/admin/dataruangan/delete/:id", middleware.CheckPrivileges("del", "9"), controller.DeleteRuanganController)
	app.Post("/admin/dataruangan/update/:id", middleware.CheckPrivileges("edit", "9"), controller.UpdateRuanganController)
	// -------------------------------------------------------------
	// ----------------------- Data Peminjaman ---------------------
	app.Get("/admin/datapeminjaman", middleware.CheckPrivileges("view", "13"), controller.DataPeminjamanController)
	app.Post("/admin/datapeminjaman/update/:id", middleware.CheckPrivileges("edit", "13"), controller.UpdatePeminjamanController)
	// -------------------------------------------------------------
	// ----------------------- Data Pengembalian -------------------
	app.Get("/admin/datapengembalian", middleware.CheckPrivileges("view", "14"), controller.DataPengembalianController)
	// -------------------------------------------------------------
	// ----------------------- Data Admin --------------------------
	app.Get("/admin/dataadmin", middleware.CheckPrivileges("view", "10"), controller.DataAdminController)
	app.Post("/admin/dataadmin/tambah", middleware.CheckPrivileges("create", "10"), controller.CreateAdminController)
	app.Post("/admin/dataadmin/delete/:id", middleware.CheckPrivileges("del", "10"), controller.DeleteAdminController)
	app.Post("/admin/dataadmin/update/:id", middleware.CheckPrivileges("edit", "10"), controller.UpdateAdminController)

	// -------------------------------------------------------------

	// app.Get("/peminjaman", controller.PeminjamanFormController)

}
