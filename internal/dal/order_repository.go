package dal

import (
	"encoding/json"
	"log/slog"
	"os"

	"hot-coffee/internal/config"
	myerrors "hot-coffee/internal/myErrors"
	"hot-coffee/internal/utils"
	"hot-coffee/models"
)

type OrderRepository interface {
	GetOrder() ([]models.Order, error)
	GetOrderID(id string) (models.Order, error)
	CreateOrder(newOrder models.Order) error
	CloseOrder(id string) error
	UpdateOrder(id string, newOrder models.Order) error
	DeleteOrder(id string) error
}

type jsonOrderRepository struct {
	filepath string
}

func NewOrderRepository(filepath string) OrderRepository {
	return &jsonOrderRepository{filepath: filepath}
}

// RETURN ERROR???
func (o *jsonOrderRepository) GetOrder() ([]models.Order, error) {
	byteValue, err := utils.ReadFile(o.filepath)
	if err != nil {
		slog.Error("Failed to open", "error", err, "file path", o.filepath)
		return []models.Order{}, myerrors.ErrFailOpenJson
	}

	var orders []models.Order
	if err := json.Unmarshal(byteValue, &orders); err != nil {
		slog.Error("Failed to unmarshal", "error", err)
		return []models.Order{}, myerrors.ErrFailUnmarshal
	}
	return orders, nil
}

// check if id exists
func (o *jsonOrderRepository) GetOrderID(id string) (models.Order, error) {
	byteValue, err := utils.ReadFile(o.filepath)
	if err != nil {
		slog.Error("Failed to open", "error", err, "file path", o.filepath)
		return models.Order{}, myerrors.ErrFailOpenJson
	}

	var orders []models.Order
	if err := json.Unmarshal(byteValue, &orders); err != nil {
		slog.Error("Failed to unmarshal", "error", err)
		return models.Order{}, myerrors.ErrFailUnmarshal
	}

	for _, order := range orders {
		if order.ID == id {
			return order, nil
		}
	}

	slog.Error("Failed to find", "error", myerrors.ErrNotFound)
	return models.Order{}, myerrors.ErrNotFound
}

func (o *jsonOrderRepository) CreateOrder(newOrder models.Order) error {
	byteValue, err := utils.ReadFile(o.filepath)
	if err != nil {
		slog.Error("Failed to open", "error", err, "file path", o.filepath)
		return myerrors.ErrFailOpenJson
	}

	var orders []models.Order
	if err := json.Unmarshal(byteValue, &orders); err != nil {
		slog.Error("Failed to unmarshal", "error", err)
		return myerrors.ErrFailUnmarshal
	}
	orders = append(orders, newOrder)

	filestring, _ := json.MarshalIndent(orders, "", "  ")
	os.WriteFile(*config.Dir+"/"+o.filepath, filestring, os.ModePerm)

	return nil
}

// check if id exists
func (o *jsonOrderRepository) CloseOrder(id string) error {
	byteValue, err := utils.ReadFile(o.filepath)
	if err != nil {
		slog.Error("Failed to open", "error", err, "file path", o.filepath)
		return myerrors.ErrFailOpenJson

	}

	var orders []models.Order
	if err := json.Unmarshal(byteValue, &orders); err != nil {
		slog.Error("Failed to unmarshal", "error", err)
		return myerrors.ErrFailUnmarshal

	}

	var isFound bool
	for i := range orders {
		if orders[i].ID == id {
			isFound = true
			if orders[i].Status == "open" {
				orders[i].Status = "closed"
			} else {
				slog.Error("Failed to close order", "error", myerrors.ErrOrderClosed)
				return myerrors.ErrOrderClosed
			}
		}
	}

	if !isFound {
		slog.Error("Failed to find", "error", myerrors.ErrNotFound)
		return myerrors.ErrNotFound
	}

	filestring, _ := json.MarshalIndent(orders, "", "  ")
	os.WriteFile(*config.Dir+"/"+o.filepath, filestring, os.ModePerm)

	return nil
}

// check if id exists
func (o *jsonOrderRepository) UpdateOrder(id string, newOrder models.Order) error {
	byteValue, err := utils.ReadFile(o.filepath)
	if err != nil {
		slog.Error("Failed to open", "error", err, "file path", o.filepath)
		return myerrors.ErrFailOpenJson
	}

	var orders []models.Order
	if err := json.Unmarshal(byteValue, &orders); err != nil {
		slog.Error("Failed to unmarshal", "error", err)
		return myerrors.ErrFailUnmarshal
	}

	var isFound bool
	for i := range orders {
		if orders[i].ID == id {
			orders[i] = newOrder
			isFound = true
		}
	}

	if !isFound {
		slog.Error("Failed to find", "error", myerrors.ErrNotFound)
		return myerrors.ErrNotFound
	}
	filestring, _ := json.MarshalIndent(orders, "", "  ")
	os.WriteFile(*config.Dir+"/"+o.filepath, filestring, os.ModePerm)

	return nil
}

// check if id exists
func (o *jsonOrderRepository) DeleteOrder(id string) error {
	byteValue, err := utils.ReadFile(o.filepath)
	if err != nil {
		slog.Error("Failed to open", "error", err, "file path", o.filepath)
		return myerrors.ErrFailOpenJson

	}

	var orders []models.Order
	if err := json.Unmarshal(byteValue, &orders); err != nil {
		slog.Error("Failed to unmarshal", "error", err)
		return myerrors.ErrFailUnmarshal
	}

	var isFound bool
	var newOrders []models.Order
	for i := range orders {
		if orders[i].ID == id {
			isFound = true
			continue
		}
		newOrders = append(newOrders, orders[i])
	}

	if !isFound {
		slog.Error("Failed to find", "error", myerrors.ErrNotFound)
		return myerrors.ErrNotFound
	}

	filestring, _ := json.MarshalIndent(newOrders, "", "  ")
	os.WriteFile(*config.Dir+"/"+o.filepath, filestring, os.ModePerm)

	return nil
}
