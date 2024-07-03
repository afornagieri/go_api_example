package usecases

import (
	"errors"

	entities "github.com/afornagieri/go_api_template/internal/domain/entities/item"
	repositories "github.com/afornagieri/go_api_template/internal/infra/repositories"
)

type UseCase interface {
	GetItems() ([]*entities.Item, error)
	GetItemsByName(name string) ([]*entities.Item, error)
	CreateItem(item *entities.Item) error
	UpdateItem(name string, item *entities.Item) error
	DeleteItem(name string) error
}

type ItemUseCase struct {
	Repo repositories.ItemRepository
}

func NewIemUseCase(repo repositories.ItemRepository) *ItemUseCase {
	return &ItemUseCase{Repo: repo}
}

func (uc *ItemUseCase) GetItems() ([]*entities.Item, error) {
	return uc.Repo.GetItems()
}

func (uc *ItemUseCase) GetItemByName(name string) (*entities.Item, error) {
	return uc.Repo.GetItemByName(name)
}

func (uc *ItemUseCase) CreateItem(item entities.Item) error {
	existingItem, err := uc.Repo.GetItemByName(item.Name)
	if err != nil && err.Error() != "item not found" {
		return err
	}
	if existingItem != nil {
		return errors.New("item already exists")
	}
	return uc.Repo.CreateItem(item)
}

func (uc *ItemUseCase) UpdateItem(name string, item entities.Item) error {
	existingItem, err := uc.Repo.GetItemByName(name)
	if err != nil {
		return err
	}
	if existingItem == nil {
		return errors.New("item does not exist")
	}
	return uc.Repo.UpdateItem(name, item)
}

func (uc *ItemUseCase) DeleteItem(name string) error {
	existingItem, err := uc.Repo.GetItemByName(name)
	if err != nil {
		return err
	}
	if existingItem != nil {
		return errors.New("item does not exist")
	}
	return uc.Repo.DeleteItem(name)
}
