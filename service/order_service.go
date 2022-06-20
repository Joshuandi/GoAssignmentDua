package service

import (
	order "GoAssignmentDua/model"
	"errors"
)

type OrderServiceInterface interface {
	OrderService(order *order.Order) (*order.Order, error)
}

type OrderService struct{}

func NewOrderService() OrderServiceInterface {
	return &OrderService{}
}

func (o *OrderService) OrderService(order *order.Order) (*order.Order, error) {
	if order.Customer_name == "" {
		return nil, errors.New("Customer Name harus di isi")
	}
	return order, nil
}
