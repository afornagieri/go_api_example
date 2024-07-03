package di

import (
	"github.com/afornagieri/go_api_template/internal/adapter/controller"
	usecases "github.com/afornagieri/go_api_template/internal/domain/usecases/item"
	"github.com/afornagieri/go_api_template/internal/infra/database"
	"github.com/afornagieri/go_api_template/internal/infra/repositories"
)

type Container struct {
	ItemController *controller.ItemController
}

func NewContainer() *Container {
	db, err := database.NewSqlCli()
	if err != nil {
		panic(err)
	}

	itemRepository := repositories.ItemRepository{DB: db}
	itemUseCase := usecases.NewIemUseCase(itemRepository)
	itemController := controller.NewItemController(*itemUseCase)

	return &Container{ItemController: itemController}
}
