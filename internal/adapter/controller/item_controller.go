package controller

import (
	"encoding/json"
	"net/http"

	entities "github.com/afornagieri/go_api_template/internal/domain/entities/item"
	usecases "github.com/afornagieri/go_api_template/internal/domain/usecases/item"
	"github.com/go-chi/chi/v5"
)

type ItemController struct {
	UseCase usecases.ItemUseCase
}

func NewItemController(useCase usecases.ItemUseCase) *ItemController {
	return &ItemController{UseCase: useCase}
}

func (ctrl *ItemController) GetItems(w http.ResponseWriter, r *http.Request) {
	items, err := ctrl.UseCase.Repo.GetItems()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if items == nil {
		items = []*entities.Item{}
	}
	json.NewEncoder(w).Encode(items)
}

func (ctrl *ItemController) GetItemByName(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	if name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	item, err := ctrl.UseCase.Repo.GetItemByName(name)
	if err != nil {
		if err.Error() == "item not found" {
			item = nil
			w.WriteHeader(http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	json.NewEncoder(w).Encode(item)
}

func (ctrl *ItemController) CreateItem(w http.ResponseWriter, r *http.Request) {
	var item *entities.Item
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = ctrl.UseCase.CreateItem(item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (ctrl *ItemController) UpdateItem(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	var item *entities.Item
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = ctrl.UseCase.UpdateItem(name, item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (ctrl *ItemController) DeleteItem(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	err := ctrl.UseCase.Repo.DeleteItem(name)
	if err != nil {
		if err.Error() != "item does not exist" {
			w.WriteHeader(http.StatusNotFound)
		}
	}
	w.WriteHeader(http.StatusNoContent)
}
