package internal

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func NewItemStorage() *ItemStorage {
	return &ItemStorage{
		items: make(map[string]Item),
	}
}

func (s *ItemStorage) Create(item Item) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := uuid.New().String()
	item.ID = id
	item.Checked = false
	s.items[id] = item

	return id, nil
}

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

func (s *ItemStorage) Update(id string, name *string, checked *bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	item, ok := s.items[id]
	if !ok {
		return fiber.ErrNotFound
	}

	if name != nil {
		item.Name = *name
	}
	if checked != nil {
		item.Checked = *checked
	}

	s.items[id] = item
	return nil
}

func (s *ItemStorage) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.items[id]; !ok {
		return fiber.ErrNotFound
	}

	delete(s.items, id)
	return nil
}
