package internal

import "github.com/gofiber/fiber/v2"

func SetupRoutes(app *fiber.App, handler *Handler) {
	app.Get("/", handler.Root)

	api := app.Group("/api")
	api.Post("/items", handler.CreateItem)
	api.Get("/items", handler.ListItems)
	api.Get("/items/:id", handler.GetItem)
	api.Patch("/items/:id", handler.UpdateItem)
	api.Delete("/items/:id", handler.DeleteItem)
}
