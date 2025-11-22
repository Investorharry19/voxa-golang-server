package routers

import (
	"github.com/Investorharry19/voxa-golang-server/controllers"
	"github.com/Investorharry19/voxa-golang-server/middlewares"
	"github.com/gofiber/fiber/v2"
)

func MessageRouter(app *fiber.App) {
	messageGroup := app.Group("/message")

	messageGroup.Post("/send/text-message", controllers.AddTextMessage)
	messageGroup.Post("/send/audio-message", controllers.SendAudioMessage)

	messageGroup.Get("/get-messages", middlewares.RequireAuth, controllers.GetAllMessages)
	messageGroup.Patch("/mark-as-read/:id", middlewares.RequireAuth, controllers.MarkAsRead)
	messageGroup.Patch("/star-message/:id", middlewares.RequireAuth, controllers.StarMessage)
	messageGroup.Delete("/delete-message/:id", middlewares.RequireAuth, controllers.DeleteOneMessage)
	messageGroup.Delete("/delete-all-messages", middlewares.RequireAuth, controllers.DeleteAllMessages)

	app.Get("/convert", controllers.HandleVideoBuffer)
	app.Post("/process", controllers.ProcessAudioMessage)
}
