package routers

import (
	"github.com/Investorharry19/voxa-golang-server/controllers"
	"github.com/Investorharry19/voxa-golang-server/middlewares"
	"github.com/gofiber/fiber/v2"
)

func UserRouter(app *fiber.App) {
	accountGroup := app.Group("/account")

	accountGroup.Get("/users", controllers.GetUsers)
	accountGroup.Post("/register", controllers.RegisterUser)
	accountGroup.Post("/login", controllers.LoginUser)
	accountGroup.Get("/current-user", middlewares.RequireAuth, controllers.GetCurrentUser)
}
