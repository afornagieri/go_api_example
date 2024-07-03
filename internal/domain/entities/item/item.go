package entities

import (
	"errors"

	"github.com/google/uuid"
)

type Item struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Price       float64   `json:"price"`
	Description string    `json:"description"`
}

func NewItem(name string, price float64, description string) (*Item, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}
	if price <= 0 {
		return nil, errors.New("price must be greater than 0")
	}
	if description == "" {
		return nil, errors.New("description is required")
	}
	return &Item{
		ID:          uuid.New(),
		Name:        name,
		Price:       price,
		Description: description,
	}, nil
}
