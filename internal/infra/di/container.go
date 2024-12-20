package di

import (
	controller "github.com/afornagieri/go_api_template/internal/adapter/controllers"
	"github.com/afornagieri/go_api_template/internal/domain/usecases"
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

	itemRepository := repositories.NewItemRepository(db)
	itemUseCase := usecases.NewItemUseCase(itemRepository)
	itemController := controller.NewItemController(itemUseCase)

	return &Container{ItemController: itemController}
}
