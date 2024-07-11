package repositories

import (
	"database/sql"
	"errors"

	entities "github.com/afornagieri/go_api_template/internal/domain/entities/item"
	database "github.com/afornagieri/go_api_template/internal/infra/database"
)

type ItemRepositoryInterface interface {
	GetItems() ([]*entities.Item, error)
	GetItemByName() (*entities.Item, error)
	CreateItem(item entities.Item) error
	UpdateItem(name string, item entities.Item) error
	DeleteItem(name string) error
}

type ItemRepository struct {
	DB *database.SqlCli
}

func NewItemRepository(db *database.SqlCli) *ItemRepository {
	return &ItemRepository{DB: db}
}

func (repo *ItemRepository) GetItems() ([]*entities.Item, error) {
	var items []*entities.Item

	rows, err := repo.DB.Conn.Query("SELECT id, name, price, description FROM items")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item entities.Item
		var id string

		err := rows.Scan(&id, &item.Name, &item.Price, &item.Description)
		if err != nil {
			return nil, err
		}

		items = append(items, &item)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (repo *ItemRepository) GetItemByName(name string) (*entities.Item, error) {
	var item entities.Item

	err := repo.DB.Conn.QueryRow("SELECT id, name, price, description FROM items WHERE name = ?", name).
		Scan(&item.ID, &item.Name, &item.Price, &item.Description)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("item not found")
		}
		return nil, err
	}

	return &item, nil
}

func (repo *ItemRepository) CreateItem(item entities.Item) error {
	itemFound, err := repo.GetItemByName(item.Name)
	if err != nil && err.Error() != "item not found" {
		return err
	}
	if itemFound != nil {
		return errors.New("item already exists")
	}

	newItem, err := entities.NewItem(item.Name, item.Price, item.Description)
	if err != nil {
		return err
	}

	_, err = repo.DB.Conn.Exec("INSERT INTO items (id, name, price, description) VALUES (?, ?, ?, ?)", newItem.ID.String(), newItem.Name, newItem.Price, newItem.Description)
	if err != nil {
		return err
	}

	return nil
}

func (repo *ItemRepository) UpdateItem(name string, item entities.Item) error {
	itemFound, err := repo.GetItemByName(name)
	if err != nil && err.Error() != "item not found" {
		return err
	}
	if itemFound == nil {
		return errors.New("item does not exist")
	}

	_, err = repo.DB.Conn.Exec("UPDATE items SET name = ?, price = ?, description = ? WHERE name = ?", item.Name, item.Price, item.Description, name)
	if err != nil {
		return err
	}

	return nil
}

func (repo *ItemRepository) DeleteItem(name string) error {
	result, err := repo.DB.Conn.Exec("DELETE FROM items WHERE name = ?", &name)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("item does not exist")
	}

	return nil
}
