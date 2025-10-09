package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/rustamyusupov/shoplist/internal"
)

func main() {
	app := fiber.New()

	app.Use(requestid.New())
	app.Use(logger.New(logger.Config{
		Format:     "${locals:requestid}: ${time} ${method} ${path} - ${status} - ${latency}\n",
		TimeFormat: "2006-01-02 15:04:05.000000",
	}))

	storage := internal.NewItemStorage()
	handler := internal.NewHandler(storage)
	internal.SetupRoutes(app, handler)

	log.Fatal(app.Listen(":8080"))
}
