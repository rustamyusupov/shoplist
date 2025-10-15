package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/template/html/v2"
	"github.com/rustamyusupov/shoplist/internal"
)

func main() {
	views := html.New("web/templates", ".tmpl")

	app := fiber.New(fiber.Config{Views: views})

	app.Use(requestid.New())
	app.Static("/", "web/static")
	app.Use(logger.New(logger.Config{
		Format:     "${locals:requestid}: ${time} ${method} ${path} - ${status} - ${latency}\n",
		TimeFormat: "2006-01-02 15:04:05.000000",
	}))

	storage := internal.NewStorage("shoplist.db")
	handler := internal.NewHandler(storage)
	internal.SetupRoutes(app, handler)

	log.Fatal(app.Listen(":8080"))
}
