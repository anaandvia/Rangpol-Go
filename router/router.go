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

	app.Get("/", controller.HomeController)
	// app.Get("/login", controller.LoginController)
	// app.Post("/login", middleware.Authenticate)
	// app.Get("/login", controller.LoginFormController) // Route for login form
	// app.Post("/login", controller.LoginUser)

	app.Get("/login", middleware.RedirectIfAuthenticated, controller.LoginFormController)
	app.Get("/detail_room", controller.RoomDetailController)
	app.Get("/peminjaman", controller.PeminjamanFormController)
	app.Post("/borang", controller.PeminjamanController)

	// Other routes
	app.Post("/login", controller.LoginPostController)
	app.Get("/logout", controller.LogoutController)
	app.Get("/admin", controller.AdminPage)
	app.Get("/test", controller.TestController)
	app.Get("/history", controller.HistoryPeminjamanController)

	app.Get("/register", controller.RegisterFormController) // Route for registration form
	app.Post("/register", controller.RegisterController)
	// app.Post("/", controller.BlogCreate)
	// app.Put("/:id", controller.BlogUpdate)
	// app.Delete("/:id", controller.BlogDelete)

}
