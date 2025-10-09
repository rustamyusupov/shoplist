package internal

import "github.com/gofiber/fiber/v2"

func SetupRoutes(app *fiber.App, handler *Handler) {
	app.Get("/", handler.Root)
	app.Post("/list", handler.CreateItem)
	app.Get("/list", handler.ListItems)
	app.Get("/list/:id", handler.GetItem)
	app.Patch("/list/:id", handler.UpdateItem)
	app.Delete("/list/:id", handler.DeleteItem)
}
