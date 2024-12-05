package validation

import (
	"hot-coffee/models"
	"log/slog"

	myerrors "hot-coffee/internal/myErrors"
)

func IsJSON(contentType string) bool {
	if contentType == "application/json" {
		return true
	}
	return false
}

func CheckOrder(newOrder models.Order) error {
	if newOrder.CustomerName == "" {
		slog.Error("Validation failed: Customer name is required")
		return myerrors.ErrNameRequired
	}
	if len(newOrder.Items) == 0 {
		slog.Error("Validation failed: Items is required")
		return myerrors.ErrItemsRequired
	} else {
		for _, item := range newOrder.Items {
			if item.ProductID == "" {
				slog.Error("Validation failed: Product ID is required")
				return myerrors.ErrIdRequired
			}
			if item.Quantity <= 0 {
				slog.Error("Validation failed: Quantity must be <=0", "quantity", item.Quantity)
				return myerrors.ErrInvalidQuantity
			}
		}
	}

	return nil
}

func CheckMenu(newMenu models.MenuItem) error {
	if newMenu.ID == "" {
		slog.Error("Validation failed: ID field is required")
		return myerrors.ErrIdRequired
	}
	if newMenu.Name == "" {
		slog.Error("Validation failed: Name field is required")
		return myerrors.ErrNameRequired
	}
	if newMenu.Description == "" {
		slog.Error("Validation failed: Description field is required")
		return myerrors.ErrDescriptionRequired
	}
	if newMenu.Price < 0 {
		slog.Error("Validation failed: Price field must be >=0")
		return myerrors.ErrPriceRequired
	}
	if len(newMenu.Ingredients) == 0 {
		slog.Error("Validation failed: Ingredients field are required")
		return myerrors.ErrIngredientsRequired
	} else {
		for _, itemIngredient := range newMenu.Ingredients {
			if itemIngredient.IngredientID == "" {
				slog.Error("Validation failed: Ingredient ID field is required")
				return myerrors.ErrIdRequired
			}
			if itemIngredient.Quantity <= 0 {
				slog.Error("Validation failed: Quantity field must be >1", "quantity", itemIngredient.Quantity)
				return myerrors.ErrInvalidQuantity
			}
		}
	}

	return nil
}

func CheckInventory(newInvent models.InventoryItem) error {
	if newInvent.IngredientID == "" {
		slog.Error("Validation failed: Ingredient field ID is required")
		return myerrors.ErrIdRequired
	}
	if newInvent.Name == "" {
		slog.Error("Validation failed: Name field is required")
		return myerrors.ErrNameRequired
	}
	if newInvent.Quantity < 0 {
		slog.Error("Validation failed: Quantity field must be >=0")
		return myerrors.ErrInvalidQuantity
	}
	if newInvent.Unit == "" {
		slog.Error("Validation failed: Unit field is required")
		return myerrors.ErrUnitRequired
	}

	return nil
}
