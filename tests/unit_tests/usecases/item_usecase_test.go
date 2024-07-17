package usecases_test

import (
	"testing"

	entities "github.com/afornagieri/go_api_template/internal/domain/entities/item"
	"github.com/afornagieri/go_api_template/internal/domain/usecases"
	"github.com/afornagieri/go_api_template/tests/unit_tests/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGetItems(t *testing.T) {
	mockRepo := mocks.NewMockItemRepository()
	usecase := usecases.NewItemUseCase(mockRepo)

	item1 := &entities.Item{Name: "Item1", Price: 10.0, Description: "Description1"}
	item2 := &entities.Item{Name: "Item2", Price: 20.0, Description: "Description2"}

	usecase.CreateItem(item1)
	usecase.CreateItem(item2)

	items, err := usecase.GetItems()
	assert.NoError(t, err)
	assert.Len(t, items, 2)
}

func TestGetItemByName(t *testing.T) {
	mockRepo := mocks.NewMockItemRepository()
	usecase := usecases.NewItemUseCase(mockRepo)

	item := &entities.Item{Name: "Item", Price: 10.0, Description: "Description"}

	usecase.CreateItem(item)

	item, err := usecase.GetItemByName(item.Name)
	assert.NoError(t, err)
	assert.NotNil(t, item)
}

func TestGetItemByName_ShouldReturnError(t *testing.T) {
	mockRepo := mocks.NewMockItemRepository()
	usecase := usecases.NewItemUseCase(mockRepo)

	item := &entities.Item{Name: "Item", Price: 10.0, Description: "Description"}

	err := usecase.CreateItem(item)
	assert.NoError(t, err)

	_, err = usecase.GetItemByName("Item1")
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "item not found")
}

func TestCreateItem(t *testing.T) {
	mockRepo := mocks.NewMockItemRepository()
	usecase := usecases.NewItemUseCase(mockRepo)

	item := &entities.Item{Name: "Item", Price: 10.0, Description: "Description"}

	err := usecase.CreateItem(item)
	assert.NoError(t, err)

	item, err = usecase.GetItemByName(item.Name)
	assert.NotNil(t, item)
	assert.NoError(t, err)
}

func TestCreateItem_ShouldReturnError(t *testing.T) {
	mockRepo := mocks.NewMockItemRepository()
	usecase := usecases.NewItemUseCase(mockRepo)

	item := &entities.Item{Name: "Item", Price: 10.0, Description: "Description"}

	err := usecase.CreateItem(item)
	assert.NoError(t, err)

	err = usecase.CreateItem(item)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "item already exists")
}

func TestUpdateItem(t *testing.T) {
	mockRepo := mocks.NewMockItemRepository()
	usecase := usecases.NewItemUseCase(mockRepo)

	oldItem := &entities.Item{Name: "Item1", Price: 10.0, Description: "Description1"}

	err := usecase.CreateItem(oldItem)
	assert.NoError(t, err)

	item, err := usecase.GetItemByName(oldItem.Name)
	assert.NotNil(t, item)
	assert.NoError(t, err)

	newItem := &entities.Item{Name: "Item2", Price: 20.0, Description: "Description2"}
	err = usecase.UpdateItem(oldItem.Name, newItem)
	assert.NoError(t, err)
}

func TestUpdateItem_ShouldReturnError(t *testing.T) {
	mockRepo := mocks.NewMockItemRepository()
	usecase := usecases.NewItemUseCase(mockRepo)

	item := &entities.Item{Name: "Item", Price: 10.0, Description: "Description"}

	err := usecase.CreateItem(item)
	assert.NoError(t, err)

	item, err = usecase.GetItemByName("Item")
	assert.NotNil(t, item)
	assert.NoError(t, err)

	err = usecase.UpdateItem("item1", item)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "item not found")
}

func TestDeleteItem(t *testing.T) {
	mockRepo := mocks.NewMockItemRepository()
	usecase := usecases.NewItemUseCase(mockRepo)

	item := &entities.Item{Name: "Item", Price: 10.0, Description: "Description"}

	err := usecase.CreateItem(item)
	assert.NoError(t, err)

	item, err = usecase.GetItemByName(item.Name)
	assert.NotNil(t, item)
	assert.NoError(t, err)

	err = usecase.DeleteItem(item.Name)
	assert.NoError(t, err)
}

func TestDeleteItem_ShouldReturnError(t *testing.T) {
	mockRepo := mocks.NewMockItemRepository()
	usecase := usecases.NewItemUseCase(mockRepo)

	item := &entities.Item{Name: "Item", Price: 10.0, Description: "Description"}

	err := usecase.CreateItem(item)
	assert.NoError(t, err)

	item, err = usecase.GetItemByName(item.Name)
	assert.NotNil(t, item)
	assert.NoError(t, err)

	err = usecase.DeleteItem("Item2")
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "item not found")
}
