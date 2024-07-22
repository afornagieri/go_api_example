package mocks

import (
	"errors"

	entities "github.com/afornagieri/go_api_template/internal/domain/entities/item"
)

type MockItemRepository struct {
	items                 map[string]*entities.Item
	shouldErrorGetItems   bool
	shouldErrorGetItem    bool
	shouldErrorCreateItem bool
	shouldErrorUpdateItem bool
	shouldErrorDeleteItem bool
}

func NewMockItemRepository() *MockItemRepository {
	return &MockItemRepository{
		items: make(map[string]*entities.Item),
	}
}

func (m *MockItemRepository) SetError(method string, shouldError bool) {
	switch method {
	case "GetItems":
		m.shouldErrorGetItems = shouldError
	case "GetItem":
		m.shouldErrorGetItem = shouldError
	case "CreateItem":
		m.shouldErrorCreateItem = shouldError
	case "UpdateItem":
		m.shouldErrorUpdateItem = shouldError
	case "DeleteItem":
		m.shouldErrorDeleteItem = shouldError
	}
}

func (m *MockItemRepository) GetItems() ([]*entities.Item, error) {
	if m.shouldErrorGetItems {
		return nil, errors.New("internal server error")
	}
	var itemList []*entities.Item
	for _, itm := range m.items {
		itemList = append(itemList, itm)
	}
	return itemList, nil
}

func (m *MockItemRepository) GetItemByName(name string) (*entities.Item, error) {
	if m.shouldErrorGetItem {
		return nil, errors.New("internal server error")
	}
	itm, exists := m.items[name]
	if !exists {
		return nil, errors.New("item not found")
	}
	return itm, nil
}

func (m *MockItemRepository) CreateItem(itm *entities.Item) error {
	if m.shouldErrorCreateItem {
		return errors.New("internal server error")
	}
	if _, exists := m.items[itm.Name]; exists {
		return errors.New("item already exists")
	}
	m.items[itm.Name] = itm
	return nil
}

func (m *MockItemRepository) UpdateItem(name string, itm *entities.Item) error {
	if m.shouldErrorUpdateItem {
		return errors.New("internal server error")
	}
	if _, exists := m.items[name]; !exists {
		return errors.New("item not found")
	}
	m.items[name] = itm
	return nil
}

func (m *MockItemRepository) DeleteItem(name string) error {
	if m.shouldErrorDeleteItem {
		return errors.New("internal server error")
	}
	if _, exists := m.items[name]; !exists {
		return errors.New("item not found")
	}
	delete(m.items, name)
	return nil
}
