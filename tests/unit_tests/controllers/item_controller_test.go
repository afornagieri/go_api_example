package controller_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/afornagieri/go_api_template/internal/adapter/controllers"
	entities "github.com/afornagieri/go_api_template/internal/domain/entities/item"
	"github.com/afornagieri/go_api_template/internal/domain/usecases"
	"github.com/afornagieri/go_api_template/tests/unit_tests/mocks"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func setupController() (*controllers.ItemController, *mocks.MockItemRepository) {
	mockRepo := mocks.NewMockItemRepository()
	useCase := usecases.NewItemUseCase(mockRepo)
	controller := controllers.NewItemController(useCase)
	return controller, mockRepo
}

func executeRequest(req *http.Request, ctrl *controllers.ItemController) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	router := chi.NewRouter()

	router.Get("/items", ctrl.GetItems)
	router.Get("/items/{name}", ctrl.GetItemByName)
	router.Post("/items", ctrl.CreateItem)
	router.Put("/items/{name}", ctrl.UpdateItem)
	router.Delete("/items/{name}", ctrl.DeleteItem)

	router.ServeHTTP(recorder, req)

	return recorder
}

func TestGetItemsController(t *testing.T) {
	ctrl, mockRepo := setupController()

	item1 := &entities.Item{Name: "item1", Price: 10.0, Description: "Description1"}
	item2 := &entities.Item{Name: "item2", Price: 10.0, Description: "Description2"}
	mockRepo.CreateItem(item1)
	mockRepo.CreateItem(item2)

	req, _ := http.NewRequest("GET", "/items", nil)
	response := executeRequest(req, ctrl)

	assert.Equal(t, http.StatusOK, response.Code)

	var items []*entities.Item
	err := json.NewDecoder(response.Body).Decode(&items)
	assert.NoError(t, err)
	assert.Len(t, items, 2)
}
