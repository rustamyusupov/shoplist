package internal

import (
	"database/sql"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	_ "github.com/mattn/go-sqlite3"
)

func NewStorage(dbPath string) *Storage {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Printf("failed to open database: %v", err)
		return nil
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	if err := db.Ping(); err != nil {
		log.Printf("failed to ping database: %v", err)
		return nil
	}

	createTableSQL := `
    CREATE TABLE IF NOT EXISTS items (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        checked BOOLEAN NOT NULL DEFAULT 0
    );`

	if _, err := db.Exec(createTableSQL); err != nil {
		log.Printf("failed to create table: %v", err)
		return nil
	}

	return &Storage{db: db}
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) Create(item Item) (string, error) {
	result, err := s.db.Exec(
		"INSERT INTO items (name, checked) VALUES (?, ?)",
		item.Name, false,
	)
	if err != nil {
		return "", err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return "", err
	}

	return strconv.FormatInt(id, 10), nil
}

func (s *Storage) List() []Item {
	rows, err := s.db.Query("SELECT id, name, checked FROM items")
	if err != nil {
		return []Item{}
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var item Item
		if err := rows.Scan(&item.ID, &item.Name, &item.Checked); err != nil {
			continue
		}
		items = append(items, item)
	}

	return items
}

func (s *Storage) Get(id string) (Item, error) {
	var item Item
	err := s.db.QueryRow(
		"SELECT id, name, checked FROM items WHERE id = ?", id,
	).Scan(&item.ID, &item.Name, &item.Checked)

	if err == sql.ErrNoRows {
		return Item{}, fiber.ErrNotFound
	}
	if err != nil {
		return Item{}, err
	}

	return item, nil
}

func (s *Storage) Update(id string, name *string, checked *bool) error {
	query := "UPDATE items SET "
	args := []interface{}{}
	updates := []string{}

	if name != nil {
		updates = append(updates, "name = ?")
		args = append(args, *name)
	}
	if checked != nil {
		updates = append(updates, "checked = ?")
		args = append(args, *checked)
	}

	if len(updates) == 0 {
		return nil
	}

	query += updates[0]
	for i := 1; i < len(updates); i++ {
		query += ", " + updates[i]
	}
	query += " WHERE id = ?"
	args = append(args, id)

	result, err := s.db.Exec(query, args...)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fiber.ErrNotFound
	}

	return nil
}

func (s *Storage) Delete(id string) error {
	result, err := s.db.Exec("DELETE FROM items WHERE id = ?", id)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fiber.ErrNotFound
	}

	return nil
}
