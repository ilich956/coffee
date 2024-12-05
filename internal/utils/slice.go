package utils

import (
	"hot-coffee/models"
)

func DeleteElement(slice []models.OrderItem, index int) []models.OrderItem {
	return append(slice[:index], slice[index+1:]...)
}

func GetInventoryID(id string, inventeryItems []models.InventoryItem) models.InventoryItem {
	for _, inventoryItem := range inventeryItems {
		if inventoryItem.IngredientID == id {
			return inventoryItem
		}
	}
	return models.InventoryItem{}
}

func DecreaseTemporaryStock(ingID string, qty float64, tempInventory []models.InventoryItem) []models.InventoryItem {
	for _, item := range tempInventory {
		if item.IngredientID == ingID {
			item.Quantity -= qty
		}
	}
	return tempInventory
}

func AggregateOrderItems(items []models.OrderItem) []models.OrderItem {
	aggregated := make(map[string]int)
	for _, item := range items {
		aggregated[item.ProductID] += item.Quantity
	}

	result := make([]models.OrderItem, 0, len(aggregated))
	for id, quantity := range aggregated {
		result = append(result, models.OrderItem{ProductID: id, Quantity: quantity})
	}

	return result
}