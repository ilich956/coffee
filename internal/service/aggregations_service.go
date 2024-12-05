package service

import (
	"hot-coffee/internal/dal"
)

type AggregationsService interface{}

type aggregationsService struct {
	menuRepo  dal.MenuRepository
	orderRepo dal.OrderRepository
}

func NewAggregationsService(menuRepo dal.MenuRepository, orderRepo dal.OrderRepository) AggregationsService {
	return &aggregationsService{
		menuRepo:  menuRepo,
		orderRepo: orderRepo,
	}
}
