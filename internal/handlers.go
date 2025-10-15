package internal

import (
	"github.com/gofiber/fiber/v2"
)

func NewHandler(storage *ItemStorage) *Handler {
	return &Handler{storage: storage}
}

func (h *Handler) Root(c *fiber.Ctx) error {
	items := h.storage.List()

	return c.Render("root", fiber.Map{
		"Title": "Shoplist",
		"Items": items,
	})
}

func (h *Handler) CreateItem(c *fiber.Ctx) error {
	var req CreateItemRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Name is required"})
	}

	id, err := h.storage.Create(Item{Name: req.Name})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create item"})
	}

	return c.Status(fiber.StatusCreated).JSON(CreateItemResponse{ID: id})
}

func (h *Handler) ListItems(c *fiber.Ctx) error {
	items := h.storage.List()

	resp := ListResponse{
		Items: make([]Item, len(items)),
	}
	for i, item := range items {
		resp.Items[i] = Item(item)
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (h *Handler) GetItem(c *fiber.Ctx) error {
	id := c.Params("id")
	item, err := h.storage.Get(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Item not found"})
	}

	return c.Status(fiber.StatusOK).JSON(GetItemResponse{Item: item})
}

func (h *Handler) UpdateItem(c *fiber.Ctx) error {
	id := c.Params("id")

	var req UpdateItemRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if req.Name == nil && req.Checked == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "At least one field (name or checked) must be provided"})
	}

	if err := h.storage.Update(id, req.Name, req.Checked); err != nil {
		if err == fiber.ErrNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Item not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not update item"})
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *Handler) DeleteItem(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.storage.Delete(id); err != nil {
		if err == fiber.ErrNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Item not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not delete item"})
	}

	return c.SendStatus(fiber.StatusOK)
}
