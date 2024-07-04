package router

import (
	"github.com/afornagieri/go_api_template/internal/adapter/controller"
	"github.com/afornagieri/go_api_template/internal/adapter/middlewares"
	"github.com/go-chi/chi/v5"
)

func NewRouter(itemController *controller.ItemController) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middlewares.Logging)

	r.Get("/items", itemController.GetItems)
	r.Get("/items/{name}", itemController.GetItemByName)
	r.Post("/items", itemController.CreateItem)
	r.Put("/items/{name}", itemController.UpdateItem)
	r.Delete("/items/{name}", itemController.DeleteItem)

	return r
}
