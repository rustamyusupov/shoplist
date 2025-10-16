package internal

import (
	"database/sql"
	"time"
)

type Storage struct {
	db *sql.DB
}
type Item struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Checked    bool      `json:"checked"`
	ModifiedAt time.Time `json:"modified_at"`
}

type Handler struct {
	storage *Storage
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
