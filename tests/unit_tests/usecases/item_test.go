package usecases_test

import (
	"testing"

	entities "github.com/afornagieri/go_api_template/internal/domain/entities/item"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type ItemRepositoryMock struct {
	mock.Mock
}

func (mockRepository *ItemRepositoryMock) GetItems() ([]*entities.Item, error) {
	args := mockRepository.Called()
	items := args.Get(0).([]entities.Item)
	result := make([]*entities.Item, len(items))
	for i, item := range items {
		result[i] = &item
	}
	return result, args.Error(1)
}

func TestGetItems(t *testing.T) {
	mockRepo := new(ItemRepositoryMock)

	items := []entities.Item{
		{ID: uuid.New(), Name: "item1", Price: 10.0, Description: "description"},
		{ID: uuid.New(), Name: "item2", Price: 20.0, Description: "description"},
	}

	mockRepo.On("GetItems").Return(items, nil)

	result, err := mockRepo.GetItems()

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, len(items), len(result))
	for i, item := range items {
		assert.Equal(t, item.ID, result[i].ID)
		assert.Equal(t, item.Name, result[i].Name)
		assert.Equal(t, item.Price, result[i].Price)
		assert.Equal(t, item.Description, result[i].Description)
	}

	mockRepo.AssertExpectations(t)
}
