package repositories

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	entities "github.com/afornagieri/go_api_template/internal/domain/entities/item"
	"github.com/afornagieri/go_api_template/internal/infra/database"
	"github.com/afornagieri/go_api_template/internal/infra/repositories"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetItemByName_ShouldReturn_Success_Given_Valid_Name(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := &repositories.ItemRepository{
		DB: &database.SqlCli{Conn: db},
	}

	expectedItem := &entities.Item{
		ID:          uuid.New(),
		Name:        "item1",
		Description: "description",
		Price:       10.0,
	}

	rows := mock.NewRows([]string{"id", "name", "price", "description"}).AddRow(expectedItem.ID.String(), expectedItem.Name, expectedItem.Price, expectedItem.Description)
	mock.ExpectQuery("SELECT id, name, price, description FROM items WHERE name = ?").WithArgs(expectedItem.Name).WillReturnRows(rows)

	item, err := repo.GetItemByName("item1")
	if err != nil {
		t.Fatalf("an error '%s' was not expected when querying the database", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NoError(t, err)
	assert.NotNil(t, item)
	assert.Equal(t, expectedItem, item)
}

func TestGetItems(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := &repositories.ItemRepository{
		DB: &database.SqlCli{Conn: db},
	}

	var expectedItems []*entities.Item

	for i := 0; i < 10; i++ {
		itemID := uuid.New()
		item := &entities.Item{
			ID:          itemID,
			Name:        "item" + fmt.Sprint(i),
			Description: "description",
		}
		expectedItems = append(expectedItems, item)
	}

	for _, item := range expectedItems {
		rows := mock.NewRows([]string{"id", "name", "price", "description"}).AddRow(item.ID.String(), item.Name, item.Price, item.Description)
		mock.ExpectQuery("SELECT id, name, price, description FROM items").WillReturnRows(rows)
	}

	items, err := repo.GetItems()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when querying the database", err)
	}

	assert.NoError(t, err)
	assert.NotNil(t, items)
}
