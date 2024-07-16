package usecases

import (
	entities "github.com/afornagieri/go_api_template/internal/domain/entities/item"
	"github.com/afornagieri/go_api_template/internal/infra/repositories"
)

type ItemUseCase struct {
	Repo repositories.ItemRepository
}

func NewItemUseCase(repo repositories.ItemRepository) *ItemUseCase {
	return &ItemUseCase{Repo: repo}
}

func (uc *ItemUseCase) GetItems() ([]*entities.Item, error) {
	return uc.Repo.GetItems()
}

func (uc *ItemUseCase) GetItemByName(name string) (*entities.Item, error) {
	return uc.Repo.GetItemByName(name)
}

func (uc *ItemUseCase) CreateItem(itm *entities.Item) error {
	return uc.Repo.CreateItem(itm)
}

func (uc *ItemUseCase) UpdateItem(name string, itm *entities.Item) error {
	return uc.Repo.UpdateItem(name, itm)
}

func (uc *ItemUseCase) DeleteItem(name string) error {
	return uc.Repo.DeleteItem(name)
}
