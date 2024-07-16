package repositories

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"

	entities "github.com/afornagieri/go_api_template/internal/domain/entities/item"
	database "github.com/afornagieri/go_api_template/internal/infra/database"
)

type ItemRepository_Impl struct {
	DB *database.SqlCli
}

func NewItemRepository(db *database.SqlCli) *ItemRepository_Impl {
	return &ItemRepository_Impl{DB: db}
}

func (repo *ItemRepository_Impl) GetItems() ([]*entities.Item, error) {
	var items []*entities.Item

	rows, err := repo.DB.Conn.Query("SELECT id, name, price, description FROM items")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch items: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item entities.Item
		var id string

		err := rows.Scan(&id, &item.Name, &item.Price, &item.Description)
		if err != nil {
			return nil, fmt.Errorf("failed to scan item row: %v", err)
		}

		item.ID, _ = uuid.Parse(id)
		items = append(items, &item)
	}

	return items, nil
}

func (repo *ItemRepository_Impl) GetItemByName(name string) (*entities.Item, error) {
	var item entities.Item

	err := repo.DB.Conn.QueryRow("SELECT id, name, price, description FROM items WHERE name = ?", name).
		Scan(&item.ID, &item.Name, &item.Price, &item.Description)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("item '%s' not found: %v", name, err)
		}
		return nil, fmt.Errorf("failed to get item by name: %v", err)
	}

	return &item, nil
}

func (repo *ItemRepository_Impl) CreateItem(item *entities.Item) error {
	tx, err := repo.DB.Conn.Begin()
	if err != nil {
		return fmt.Errorf("could not begin transaction: %v", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	newItem, err := entities.NewItem(item.Name, item.Price, item.Description)
	if err != nil {
		return fmt.Errorf("failed to create new item: %v", err)
	}

	_, err = tx.Exec("INSERT INTO items (id, name, price, description) VALUES (?, ?, ?, ?)", newItem.ID.String(), newItem.Name, newItem.Price, newItem.Description)
	if err != nil {
		return fmt.Errorf("failed to insert item: %v", err)
	}

	err = tx.Commit()
	return nil
}

func (repo *ItemRepository_Impl) UpdateItem(name string, item *entities.Item) error {
	tx, err := repo.DB.Conn.Begin()
	if err != nil {
		return fmt.Errorf("could not begin transaction: %v", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	_, err = repo.getItemByNameInTx(tx, name)
	if err != nil {
		return fmt.Errorf("failed to get item '%s': %v", name, err)
	}

	_, err = tx.Exec("UPDATE items SET name = ?, price = ?, description = ? WHERE name = ?", item.Name, item.Price, item.Description, name)
	if err != nil {
		return fmt.Errorf("failed to update item: %v", err)
	}

	err = tx.Commit()
	return nil
}

func (repo *ItemRepository_Impl) DeleteItem(name string) error {
	tx, err := repo.DB.Conn.Begin()
	if err != nil {
		return fmt.Errorf("could not begin transaction: %v", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	_, err = tx.Exec("DELETE FROM items WHERE name = ?", name)
	if err != nil {
		return fmt.Errorf("failed to delete item: %v", err)
	}

	err = tx.Commit()
	return nil
}

func (repo *ItemRepository_Impl) getItemByNameInTx(tx *sql.Tx, name string) (*entities.Item, error) {
	var item entities.Item

	err := tx.QueryRow("SELECT id, name, price, description FROM items WHERE name = ?", name).
		Scan(&item.ID, &item.Name, &item.Price, &item.Description)

	if err != nil {
		return nil, fmt.Errorf("failed to get item '%s' in transaction: %v", name, err)
	}

	return &item, nil
}
