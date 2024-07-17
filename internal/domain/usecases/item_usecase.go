package usecases

import (
	entities "github.com/afornagieri/go_api_template/internal/domain/entities/item"
	"github.com/afornagieri/go_api_template/internal/infra/repositories"
)

type ItemUseCase_Impl struct {
	Repo repositories.ItemRepository
}

func NewItemUseCase(repo repositories.ItemRepository) *ItemUseCase_Impl {
	return &ItemUseCase_Impl{Repo: repo}
}

func (uc *ItemUseCase_Impl) GetItems() ([]*entities.Item, error) {
	return uc.Repo.GetItems()
}

func (uc *ItemUseCase_Impl) GetItemByName(name string) (*entities.Item, error) {
	return uc.Repo.GetItemByName(name)
}

func (uc *ItemUseCase_Impl) CreateItem(itm *entities.Item) error {
	return uc.Repo.CreateItem(itm)
}

func (uc *ItemUseCase_Impl) UpdateItem(name string, itm *entities.Item) error {
	return uc.Repo.UpdateItem(name, itm)
}

func (uc *ItemUseCase_Impl) DeleteItem(name string) error {
	return uc.Repo.DeleteItem(name)
}
