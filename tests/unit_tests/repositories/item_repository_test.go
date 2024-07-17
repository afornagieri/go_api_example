package repositories_test

import (
	"database/sql"
	"errors"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	entities "github.com/afornagieri/go_api_template/internal/domain/entities/item"
	"github.com/afornagieri/go_api_template/internal/infra/database"
	"github.com/afornagieri/go_api_template/internal/infra/repositories"
)

func TestItemRepository_GetItems(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlCli := &database.SqlCli{Conn: db}
	repo := repositories.NewItemRepository(sqlCli)

	t.Run("GetItems should return items successfully", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "price", "description"}).
			AddRow(uuid.New().String(), "Item1", 100.0, "Description1").
			AddRow(uuid.New().String(), "Item2", 150.0, "Description2")
		mock.ExpectQuery("SELECT id, name, price, description FROM items").
			WillReturnRows(rows)

		items, err := repo.GetItems()
		assert.NoError(t, err)
		assert.Len(t, items, 2)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GetItems should handle database error", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, name, price, description FROM items").
			WillReturnError(errors.New("database error"))

		items, err := repo.GetItems()
		assert.Error(t, err)
		assert.Nil(t, items)
		assert.EqualError(t, err, "failed to fetch items: database error")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GetItems should handle scanning error", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, name, price, description FROM items").
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "price", "description"}).
				AddRow(uuid.New().String(), "Item1", 100.0, "Description1").
				AddRow(uuid.New().String(), "Item2", "invalid_price", "Description2"))

		items, err := repo.GetItems()
		assert.Error(t, err)
		assert.Nil(t, items)
		assert.Contains(t, err.Error(), "failed to scan item row:")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestItemRepository_GetItemByName(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlCli := &database.SqlCli{Conn: db}
	repo := repositories.NewItemRepository(sqlCli)

	t.Run("GetItemByName should return item successfully", func(t *testing.T) {
		itemName := "Item1"
		rows := sqlmock.NewRows([]string{"id", "name", "price", "description"}).
			AddRow(uuid.New().String(), itemName, 100.0, "Description1")
		mock.ExpectQuery("SELECT id, name, price, description FROM items WHERE name = ?").
			WithArgs(itemName).
			WillReturnRows(rows)

		item, err := repo.GetItemByName(itemName)
		assert.NoError(t, err)
		assert.NotNil(t, item)
		assert.Equal(t, itemName, item.Name)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GetItemByName should handle item not found", func(t *testing.T) {
		itemName := "NonExistingItem"
		mock.ExpectQuery("SELECT id, name, price, description FROM items WHERE name = ?").
			WithArgs(itemName).
			WillReturnError(sql.ErrNoRows)

		item, err := repo.GetItemByName(itemName)
		assert.Error(t, err)
		assert.Nil(t, item)
		assert.EqualError(t, err, fmt.Sprintf("item '%s' not found: %v", itemName, sql.ErrNoRows))
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GetItemByName should handle database error", func(t *testing.T) {
		itemName := "Item1"
		mock.ExpectQuery("SELECT id, name, price, description FROM items WHERE name = ?").
			WithArgs(itemName).
			WillReturnError(errors.New("database error"))

		item, err := repo.GetItemByName(itemName)
		assert.Error(t, err)
		assert.Nil(t, item)
		assert.EqualError(t, err, fmt.Sprintf("failed to get item by name: %v", errors.New("database error")))
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestItemRepository_CreateItem(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlCli := &database.SqlCli{Conn: db}
	repo := repositories.NewItemRepository(sqlCli)

	t.Run("CreateItem should create item successfully", func(t *testing.T) {
		item := &entities.Item{
			Name:        "NewItem",
			Price:       200.0,
			Description: "Description for NewItem",
		}

		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO items").
			WithArgs(sqlmock.AnyArg(), item.Name, item.Price, item.Description).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.CreateItem(item)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("CreateItem should handle database begin transaction error", func(t *testing.T) {
		item := &entities.Item{
			Name:        "item",
			Price:       200.0,
			Description: "Description for item",
		}

		mock.ExpectBegin().WillReturnError(errors.New("could not begin transaction:"))

		err := repo.CreateItem(item)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "could not begin transaction")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("CreateItem should handle missing item name", func(t *testing.T) {
		item := &entities.Item{
			Name:        "",
			Price:       200.0,
			Description: "Description for NewItem",
		}
		mock.ExpectBegin()
		err := repo.CreateItem(item)
		assert.Error(t, err)
		assert.EqualError(t, err, "failed to create new item: name is required")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("CreateItem should handle missing item price is invalid", func(t *testing.T) {
		item := &entities.Item{
			Name:        "item",
			Price:       -1,
			Description: "Description for NewItem",
		}
		mock.ExpectBegin()
		err := repo.CreateItem(item)
		assert.Error(t, err)
		assert.EqualError(t, err, "failed to create new item: price must be greater than 0")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("CreateItem should handle missing item description", func(t *testing.T) {
		item := &entities.Item{
			Name:        "item",
			Price:       200.0,
			Description: "",
		}
		mock.ExpectBegin()
		err := repo.CreateItem(item)
		assert.Error(t, err)
		assert.EqualError(t, err, "failed to create new item: description is required")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("CreateItem should handle failed to insert item error", func(t *testing.T) {
		item := &entities.Item{
			Name:        "NewItem",
			Price:       200.0,
			Description: "Description for NewItem",
		}

		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO items").
			WithArgs(sqlmock.AnyArg(), item.Name, item.Price, item.Description).
			WillReturnError(errors.New("failed to insert item:"))

		err := repo.CreateItem(item)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to insert item:")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestItemRepository_UpdateItem(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlCli := &database.SqlCli{Conn: db}
	repo := repositories.NewItemRepository(sqlCli)

	t.Run("UpdateItem should update an existing item successfully", func(t *testing.T) {
		existingItemName := "ExistingItem"
		item := &entities.Item{
			Name:        "UpdatedItem",
			Price:       300.0,
			Description: "Updated Description",
		}

		mock.ExpectBegin()
		mock.ExpectQuery("SELECT id, name, price, description FROM items WHERE name = ?").
			WithArgs(existingItemName).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "price", "description"}).
				AddRow(uuid.New().String(), existingItemName, 200.0, "Original Description"))
		mock.ExpectExec("UPDATE items").
			WithArgs(item.Name, item.Price, item.Description, existingItemName).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.UpdateItem(existingItemName, item)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("UpdateItem should handle database begin transaction error", func(t *testing.T) {
		item := &entities.Item{
			Name:        "item",
			Price:       200.0,
			Description: "Description for item",
		}

		mock.ExpectBegin().WillReturnError(errors.New("could not begin transaction:"))

		err := repo.UpdateItem("item", item)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "could not begin transaction")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("UpdateItem should handle item not found", func(t *testing.T) {
		nonExistingItemName := "NonExistingItem"
		item := &entities.Item{
			Name:        "UpdatedItem",
			Price:       300.0,
			Description: "Updated Description",
		}

		mock.ExpectBegin()
		mock.ExpectQuery("SELECT id, name, price, description FROM items WHERE name = ?").
			WithArgs(nonExistingItemName).
			WillReturnError(sql.ErrNoRows)
		mock.ExpectRollback()

		err := repo.UpdateItem(nonExistingItemName, item)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("UpdateItem should handle database error on update", func(t *testing.T) {
		itemName := "ExistingItem"
		item := &entities.Item{
			Name:        "UpdatedItem",
			Price:       300.0,
			Description: "Updated Description",
		}

		mock.ExpectBegin()
		mock.ExpectQuery("SELECT id, name, price, description FROM items WHERE name = ?").
			WithArgs(itemName).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "price", "description"}).
				AddRow(uuid.New().String(), itemName, 200.0, "Original Description"))
		mock.ExpectExec("UPDATE items").
			WithArgs(item.Name, item.Price, item.Description, itemName).
			WillReturnError(errors.New("database error"))
		mock.ExpectRollback()

		err := repo.UpdateItem(itemName, item)
		assert.Error(t, err)
		assert.EqualError(t, err, fmt.Sprintf("failed to update item: %v", errors.New("database error")))
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestItemRepository_DeleteItem(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlCli := &database.SqlCli{Conn: db}
	repo := repositories.NewItemRepository(sqlCli)

	t.Run("DeleteItem should delete item successfully", func(t *testing.T) {
		itemName := "ItemToDelete"

		mock.ExpectBegin()
		mock.ExpectExec("DELETE FROM items WHERE name = ?").
			WithArgs(itemName).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.DeleteItem(itemName)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DeleteItem should handle database begin transaction error", func(t *testing.T) {
		mock.ExpectBegin().WillReturnError(errors.New("could not begin transaction:"))

		err := repo.DeleteItem("item")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "could not begin transaction")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DeleteItem should handle database error on delete", func(t *testing.T) {
		itemName := "ItemToDelete"

		mock.ExpectBegin()
		mock.ExpectExec("DELETE FROM items WHERE name = ?").
			WithArgs(itemName).
			WillReturnError(errors.New("database error"))
		mock.ExpectRollback()

		err := repo.DeleteItem(itemName)
		assert.Error(t, err)
		assert.EqualError(t, err, fmt.Sprintf("failed to delete item: %v", errors.New("database error")))
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
