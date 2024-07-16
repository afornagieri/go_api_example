package repositories

import (
	entities "github.com/afornagieri/go_api_template/internal/domain/entities/item"
)

type ItemRepository interface {
	GetItems() ([]*entities.Item, error)
	GetItemByName(name string) (*entities.Item, error)
	CreateItem(item *entities.Item) error
	UpdateItem(name string, item *entities.Item) error
	DeleteItem(name string) error
}
