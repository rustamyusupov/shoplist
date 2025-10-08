package main

import (
	"log"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/google/uuid"
)

type Item struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Checked bool   `json:"checked"`
}

type ItemStorage struct {
	mu    sync.Mutex
	items map[string]Item
}

type (
	CreateItemRequest struct {
		Name string `json:"name" validate:"required"`
	}

	CreateItemResponse struct {
		ID string `json:"id"`
	}
)

func (s *ItemStorage) Create(item Item) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := uuid.New().String()
	item.ID = id
	item.Checked = false
	s.items[id] = item

	return id, nil
}

type (
	ListResponse struct {
		Items []Item `json:"items"`
	}

	GetItemResponse struct {
		Item Item `json:"item"`
	}
)

func (s *ItemStorage) List() []Item {
	s.mu.Lock()
	defer s.mu.Unlock()

	items := make([]Item, 0, len(s.items))
	for _, item := range s.items {
		items = append(items, item)
	}

	return items
}

func (s *ItemStorage) Get(id string) (Item, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	item, ok := s.items[id]
	if !ok {
		return Item{}, fiber.ErrNotFound
	}

	return item, nil
}

func main() {
	app := fiber.New()

	storage := &ItemStorage{
		items: make(map[string]Item),
	}

	app.Use(requestid.New())
	app.Use(logger.New(logger.Config{
		Format:     "${locals:requestid}: ${time} ${method} ${path} - ${status} - ${latency}\n",
		TimeFormat: "2006-01-02 15:04:05.000000",
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).SendString("Welcome to Shoplist")
	})

	app.Post("/list", func(c *fiber.Ctx) error {
		var req CreateItemRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
		}

		if req.Name == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Name is required"})
		}

		id, err := storage.Create(Item{Name: req.Name})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create item"})
		}

		return c.Status(fiber.StatusOK).JSON(CreateItemResponse{ID: id})
	})

	app.Get("/list", func(c *fiber.Ctx) error {
		items := storage.List()

		resp := ListResponse{
			Items: make([]Item, len(items)),
		}
		for i, item := range items {
			resp.Items[i] = Item(item)
		}

		return c.Status(fiber.StatusOK).JSON(resp)
	})

	app.Get("/list/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		item, err := storage.Get(id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Item not found"})
		}

		return c.Status(fiber.StatusOK).JSON(GetItemResponse{Item: item})
	})

	// app.Patch("/list/:id", )
	// app.Delete("/list/:id", )

	log.Fatal(app.Listen(":8080"))
}
