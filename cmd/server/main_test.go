package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/rustamyusupov/shoplist/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestApp() (*fiber.App, *internal.ItemStorage) {
	app := fiber.New()
	storage := internal.NewItemStorage()
	handler := internal.NewHandler(storage)
	internal.SetupRoutes(app, handler)
	return app, storage
}

func TestRootEndpoint(t *testing.T) {
	app, _ := setupTestApp()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestCreateItem(t *testing.T) {
	app, _ := setupTestApp()

	reqBody := internal.CreateItemRequest{Name: "Test Item"}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/list", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

	var response internal.CreateItemResponse
	json.NewDecoder(resp.Body).Decode(&response)
	assert.NotEmpty(t, response.ID)
}

func TestCreateItemEmptyName(t *testing.T) {
	app, _ := setupTestApp()

	reqBody := internal.CreateItemRequest{Name: ""}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/list", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestCreateItemInvalidBody(t *testing.T) {
	app, _ := setupTestApp()

	req := httptest.NewRequest(http.MethodPost, "/list", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestListItems(t *testing.T) {
	app, storage := setupTestApp()

	storage.Create(internal.Item{Name: "Item 1"})
	storage.Create(internal.Item{Name: "Item 2"})

	req := httptest.NewRequest(http.MethodGet, "/list", nil)
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response internal.ListResponse
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, 2, len(response.Items))
}

func TestListItemsEmpty(t *testing.T) {
	app, _ := setupTestApp()

	req := httptest.NewRequest(http.MethodGet, "/list", nil)
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response internal.ListResponse
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, 0, len(response.Items))
}

func TestGetItem(t *testing.T) {
	app, storage := setupTestApp()

	id, _ := storage.Create(internal.Item{Name: "Test Item"})

	req := httptest.NewRequest(http.MethodGet, "/list/"+id, nil)
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response internal.GetItemResponse
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "Test Item", response.Item.Name)
	assert.Equal(t, id, response.Item.ID)
}

func TestGetItemNotFound(t *testing.T) {
	app, _ := setupTestApp()

	req := httptest.NewRequest(http.MethodGet, "/list/nonexistent", nil)
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestUpdateItemName(t *testing.T) {
	app, storage := setupTestApp()

	id, _ := storage.Create(internal.Item{Name: "Old Name"})

	newName := "New Name"
	reqBody := internal.UpdateItemRequest{Name: &newName}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPatch, "/list/"+id, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	item, _ := storage.Get(id)
	assert.Equal(t, "New Name", item.Name)
}

func TestUpdateItemChecked(t *testing.T) {
	app, storage := setupTestApp()

	id, _ := storage.Create(internal.Item{Name: "Test Item"})

	checked := true
	reqBody := internal.UpdateItemRequest{Checked: &checked}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPatch, "/list/"+id, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	item, _ := storage.Get(id)
	assert.True(t, item.Checked)
}

func TestUpdateItemNoFields(t *testing.T) {
	app, storage := setupTestApp()

	id, _ := storage.Create(internal.Item{Name: "Test Item"})

	reqBody := internal.UpdateItemRequest{}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPatch, "/list/"+id, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestUpdateItemNotFound(t *testing.T) {
	app, _ := setupTestApp()

	newName := "New Name"
	reqBody := internal.UpdateItemRequest{Name: &newName}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPatch, "/list/nonexistent", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestDeleteItem(t *testing.T) {
	app, storage := setupTestApp()

	id, _ := storage.Create(internal.Item{Name: "Test Item"})

	req := httptest.NewRequest(http.MethodDelete, "/list/"+id, nil)
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	_, err = storage.Get(id)
	assert.Equal(t, fiber.ErrNotFound, err)
}

func TestDeleteItemNotFound(t *testing.T) {
	app, _ := setupTestApp()

	req := httptest.NewRequest(http.MethodDelete, "/list/nonexistent", nil)
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}
