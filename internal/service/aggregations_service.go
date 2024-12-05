package service

import (
	"encoding/json"
	"fmt"
	"hot-coffee/internal/dal"
	"hot-coffee/models"
	"log/slog"
	"sort"

	myerrors "hot-coffee/internal/myErrors"
)

type AggregationsService interface {
	ServiceGetTotal() ([]byte, error)
	ServiceGetPopular() ([]byte, error)
}

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

func (a *aggregationsService) ServiceGetTotal() ([]byte, error) {
	orders, err := a.orderRepo.GetOrder()
	if err != nil {
		return nil, err
	}

	productQuantities := make(map[string]int)
	for _, order := range orders {
		for _, item := range order.Items {
			productQuantities[item.ProductID] += item.Quantity
		}
	}

	menuItems, err := a.menuRepo.GetMenu()
	if err != nil {
		return nil, err
	}

	totalSaleCount := 0.0
	for _, menuItem := range menuItems {
		if q, ok := productQuantities[menuItem.ID]; ok {
			totalSaleCount += menuItem.Price * float64(q)
		}
	}

	totalSale := models.TotalSales{
		TotalSale: totalSaleCount,
	}
	jsonFile, err := json.MarshalIndent(totalSale, "", "  ")
	if err != nil {
		slog.Error("Failed to marshal", "error", err)
		return nil, myerrors.ErrFailMarshal
	}

	return jsonFile, nil
}

func (a *aggregationsService) ServiceGetPopular() ([]byte, error) {
	orders, err := a.orderRepo.GetOrder()
	if err != nil {
		return nil, err
	}

	productQuantities := make(map[string]int)
	for _, order := range orders {
		for _, item := range order.Items {
			productQuantities[item.ProductID] += item.Quantity
		}
	}

	var popularItems []string
	for productID := range productQuantities {
		popularItems = append(popularItems, productID)
	}

	sort.SliceStable(popularItems, func(i, j int) bool {
		return productQuantities[popularItems[i]] > productQuantities[popularItems[j]]
	})

	var popItemsModels []models.PopularItem
	c := 0
	for _, productID := range popularItems {
		if c == 3 {
			break
		}
		popItem := models.PopularItem{
			ID:       productID,
			Quantity: productQuantities[productID],
		}
		popItemsModels = append(popItemsModels, popItem)
		c++
	}

	if len(popItemsModels) == 0 {
		return nil, myerrors.ErrNoItems
	}
	fmt.Println(popItemsModels)
	jsonFile, err := json.MarshalIndent(popItemsModels, "", "  ")
	if err != nil {
		slog.Error("Failed to marshal", "error", err)
		return nil, myerrors.ErrFailMarshal
	}

	return jsonFile, nil
}
