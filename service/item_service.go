package service

import (
	item "GoAssignmentDua/model"
	"errors"
)

type ItemServiceInterface interface {
	ItemService(item *item.Item) (*item.Item, error)
}

type ItemService struct{}

func NewItemService() ItemServiceInterface {
	return &ItemService{}
}

func (i *ItemService) ItemService(item *item.Item) (*item.Item, error) {
	if item.Quantity == 0 {
		return nil, errors.New("Input Quantity")
	}
	if item.Item_code == "" {
		return nil, errors.New("Input Item Code")
	}
	if item.Description == "" {
		return nil, errors.New("Input Item Description")
	}
	return item, nil
}
