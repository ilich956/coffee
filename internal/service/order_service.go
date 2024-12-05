package service

import (
	"encoding/json"
	"hot-coffee/internal/dal"
	"hot-coffee/internal/utils"
	"hot-coffee/internal/utils/uuid"
	"hot-coffee/internal/utils/validation"
	"hot-coffee/models"
	"log/slog"
	"time"

	myerrors "hot-coffee/internal/myErrors"
)

type OrderService interface {
	ServiceGetOrder() ([]byte, error)
	ServiceGetOrderID(id string) ([]byte, error)
	ServicePostOrder(newOrderByte []byte) error
	ServicePostOrderClose(id string) error
	ServicePutOrderID(id string, newOrderByte []byte) error
	ServiceDeleteOrder(id string) error
}

type orderService struct {
	orderRepo dal.OrderRepository
	menuRepo  dal.MenuRepository
	inventory dal.InventoryRepository
}

func NewOrderService(orderRepo dal.OrderRepository, menuRepo dal.MenuRepository, inventoryRepo dal.InventoryRepository) OrderService {
	return &orderService{
		orderRepo: orderRepo,
		menuRepo:  menuRepo,
		inventory: inventoryRepo,
	}
}

// CHECK JSON STRUCTURE DOES IT HAVE FIELD IN STRUCT
// SET CREATED TIME AND STATUS

// Retrieve all orders.
func (s *orderService) ServiceGetOrder() ([]byte, error) {
	ordersStruct, err := s.orderRepo.GetOrder()
	if err != nil {
		return nil, err
	}

	jsonFile, err := json.MarshalIndent(ordersStruct, "", "  ")
	if err != nil {
		return nil, err
	}
	return jsonFile, nil
}

// Retrieve a specific order by ID.
func (s *orderService) ServiceGetOrderID(id string) ([]byte, error) {
	order, err := s.orderRepo.GetOrderID(id)
	if err != nil {
		return nil, err
	}

	jsonFile, err := json.MarshalIndent(order, "", "  ")
	if err != nil {
		return nil, err
	}
	return jsonFile, nil
}

// Create a new order.
func (s *orderService) ServicePostOrder(newOrderByte []byte) error {
	var newOrder models.Order
	if err := json.Unmarshal(newOrderByte, &newOrder); err != nil {
		slog.Error("Failed to unmarshal", "error", err)
		return myerrors.ErrFailUnmarshal
	}

	if err := validation.CheckOrder(newOrder); err != nil {
		return err
	}

	items := utils.AggregateOrderItems(newOrder.Items)

	menu, err := s.menuRepo.GetMenu()
	if err != nil {
		return err
	}

	amountOfExists := 0
	var isNotExistIndex []int

	for i, item := range items {
		product := item.ProductID
		isFound := false

		for _, menuItem := range menu {
			if product == menuItem.ID {
				amountOfExists++
				isFound = true
			}
		}

		if isFound == false {
			isNotExistIndex = append(isNotExistIndex, i)
		}
	}

	if amountOfExists != len(items) {
		for i := len(isNotExistIndex) - 1; i >= 0; i-- {
			items = utils.DeleteElement(items, isNotExistIndex[i])
		}
		slog.Warn("Some items were deleted from order because they were not listed in menu")
	}

	if len(items) == 0 {
		return myerrors.ErrAbsentItem
	}

	tempInventory, err := s.inventory.GetInventory()
	if err != nil {
		return err // ok
	}

	for i := 0; i < len(items); i++ {
		item := items[i]
		menuItem, err := s.menuRepo.GetMenuID(item.ProductID)
		if err != nil {
			return err // ok
		}

		hasEnoughIngredients := true
		requiredIngredients := make(map[string]float64)

		for _, ingredient := range menuItem.Ingredients {
			inventory := utils.GetInventoryID(ingredient.IngredientID, tempInventory)

			requiredQty := float64(item.Quantity) * ingredient.Quantity
			if requiredQty > inventory.Quantity {
				hasEnoughIngredients = false
				break
			}

			requiredIngredients[ingredient.IngredientID] += requiredQty
		}

		if !hasEnoughIngredients {
			items = append(items[:i], items[i+1:]...)
			i--
		} else {
			for ingID, qty := range requiredIngredients {
				tempInventory = utils.DecreaseTemporaryStock(ingID, qty, tempInventory)
			}
		}
	}

	if len(items) == 0 {
		return myerrors.ErrEmptyOrder // ok
	}

	newOrder.Items = items
	newOrder.ID = uuid.RandStringBytesMask()
	slog.Info(newOrder.ID)
	newOrder.Status = "open"
	newOrder.CreatedAt = time.Now().Format(time.RFC3339)

	err = s.orderRepo.CreateOrder(newOrder)
	if err != nil {
		return err // ok
	}
	return nil
}

// Close an order.
func (s *orderService) ServicePostOrderClose(id string) error {
	tempInventory, err := s.inventory.GetInventory()
	if err != nil {
		return err // ok
	}

	orderr, err := s.orderRepo.GetOrderID(id)
	items := orderr.Items

	for i := 0; i < len(items); i++ {
		item := items[i]
		menuItem, err := s.menuRepo.GetMenuID(item.ProductID)
		if err != nil {
			return err // ok
		}

		hasEnoughIngredients := true
		requiredIngredients := make(map[string]float64)

		for _, ingredient := range menuItem.Ingredients {
			inventory := utils.GetInventoryID(ingredient.IngredientID, tempInventory)

			requiredQty := float64(item.Quantity) * ingredient.Quantity
			if requiredQty > inventory.Quantity {
				hasEnoughIngredients = false
				break
			}

			requiredIngredients[ingredient.IngredientID] += requiredQty
		}

		if !hasEnoughIngredients {
			return myerrors.ErrNotEnoughIngridients
		} else {

			for ingID, qty := range requiredIngredients {
				tempInventory = utils.DecreaseTemporaryStock(ingID, qty, tempInventory)
			}
			for _, invitem := range tempInventory {
				err = s.inventory.UpdateInventory(invitem.IngredientID, invitem)
			}
		}
	}
	err = s.orderRepo.CloseOrder(id)
	if err != nil {
		return err
	}

	return nil
}

// Update an existing order.
func (s *orderService) ServicePutOrderID(id string, newOrderByte []byte) error {
	var newOrder models.Order
	if err := json.Unmarshal(newOrderByte, &newOrder); err != nil {
		slog.Error("Failed to unmarshal", "error", err)
		return err
	}

	if err := validation.CheckOrder(newOrder); err != nil {
		return err
	}

	checkOrder, err := s.orderRepo.GetOrderID(id)
	if err != nil {
		return err
	}
	if checkOrder.Status == "closed" {
		slog.Error("Failed to update: status is closed", "id", id)
		return myerrors.ErrOrderClosed
	}

	newOrder.ID = checkOrder.ID
	newOrder.Status = checkOrder.Status
	newOrder.CreatedAt = checkOrder.CreatedAt
	err = s.orderRepo.UpdateOrder(id, newOrder)
	if err != nil {
		return err
	}
	return nil
}

// Delete an order.
func (s *orderService) ServiceDeleteOrder(id string) error {
	err := s.orderRepo.DeleteOrder(id)
	if err != nil {
		return err
	}
	return nil
}
