package mocks

import (
	"errors"

	entities "github.com/afornagieri/go_api_template/internal/domain/entities/item"
)

type MockItemRepository struct {
	items map[string]*entities.Item
}

func NewMockItemRepository() *MockItemRepository {
	return &MockItemRepository{
		items: make(map[string]*entities.Item),
	}
}

func (m *MockItemRepository) GetItems() ([]*entities.Item, error) {
	var itemList []*entities.Item
	for _, itm := range m.items {
		itemList = append(itemList, itm)
	}
	return itemList, nil
}

func (m *MockItemRepository) GetItemByName(name string) (*entities.Item, error) {
	itm, exists := m.items[name]
	if !exists {
		return nil, errors.New("item not found")
	}
	return itm, nil
}

func (m *MockItemRepository) CreateItem(itm *entities.Item) error {
	if _, exists := m.items[itm.Name]; exists {
		return errors.New("item already exists")
	}
	m.items[itm.Name] = itm
	return nil
}

func (m *MockItemRepository) UpdateItem(name string, itm *entities.Item) error {
	if _, exists := m.items[name]; !exists {
		return errors.New("item not found")
	}
	m.items[name] = itm
	return nil
}

func (m *MockItemRepository) DeleteItem(name string) error {
	if _, exists := m.items[name]; !exists {
		return errors.New("item not found")
	}
	delete(m.items, name)
	return nil
}
