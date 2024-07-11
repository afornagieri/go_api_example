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

func TestGetItems(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := &repositories.ItemRepository{
		DB: &database.SqlCli{Conn: db},
	}

	t.Run("GetItems should return one or more items", func(t *testing.T) {
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

		rows := mock.NewRows([]string{"id", "name", "price", "description"})
		for _, item := range expectedItems {
			rows.AddRow(item.ID.String(), item.Name, item.Price, item.Description)
		}
		mock.ExpectQuery("SELECT id, name, price, description FROM items").WillReturnRows(rows)

		items, err := repo.GetItems()
		assert.NoError(t, err)
		assert.NotNil(t, items)
		assert.Equal(t, len(expectedItems), len(items))

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("GetItems should return error on query failure", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, name, price, description FROM items").WillReturnError(fmt.Errorf("query failure"))

		items, err := repo.GetItems()
		assert.Error(t, err)
		assert.Nil(t, items)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("GetItems should return error on scan lines", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, name, price, description FROM items").WillReturnError(fmt.Errorf("scan failure"))

		items, err := repo.GetItems()
		assert.Error(t, err)
		assert.Nil(t, items)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("GetItems should return error on rows error", func(t *testing.T) {
		rows := mock.NewRows([]string{"id", "name", "price", "description"}).
			AddRow(uuid.New().String(), "item1", 10.0, "description1").
			RowError(0, fmt.Errorf("row error"))

		mock.ExpectQuery("SELECT id, name, price, description FROM items").WillReturnRows(rows)

		items, err := repo.GetItems()
		assert.Error(t, err)
		assert.Nil(t, items)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}
