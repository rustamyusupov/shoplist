package internal

import "sync"

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

type (
	ListResponse struct {
		Items []Item `json:"items"`
	}

	GetItemResponse struct {
		Item Item `json:"item"`
	}
)

type (
	UpdateItemRequest struct {
		Name    *string `json:"name,omitempty"`
		Checked *bool   `json:"checked,omitempty"`
	}
)
