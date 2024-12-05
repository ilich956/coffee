package dal

import (
	"encoding/json"
	"hot-coffee/internal/config"
	"hot-coffee/internal/utils"
	"hot-coffee/models"
	"log/slog"
	"os"

	myerrors "hot-coffee/internal/myErrors"
)

type InventoryRepository interface {
	GetInventory() ([]models.InventoryItem, error)
	GetInventoryID(id string) (models.InventoryItem, error)
	CreateInventory(newInventoryItem models.InventoryItem) error
	UpdateInventory(id string, newInvItem models.InventoryItem) error
	DeleteInventory(id string) error
}

type jsonInventoryRepository struct {
	filepath string
}

func NewInventoryRepository(filepath string) InventoryRepository {
	return &jsonInventoryRepository{filepath: filepath}
}

func (o *jsonInventoryRepository) GetInventory() ([]models.InventoryItem, error) {
	byteValue, err := utils.ReadFile(o.filepath)
	if err != nil {
		slog.Error("Failed to open", "error", err, "file path", o.filepath)
		return []models.InventoryItem{}, myerrors.ErrFailOpenJson
	}

	var inventeryItems []models.InventoryItem

	if err := json.Unmarshal(byteValue, &inventeryItems); err != nil {
		slog.Error("Failed to unmarshal", "error", err)
		return []models.InventoryItem{}, myerrors.ErrFailUnmarshal
	}

	return inventeryItems, nil
}

func (o *jsonInventoryRepository) GetInventoryID(id string) (models.InventoryItem, error) {
	byteValue, err := utils.ReadFile(o.filepath)
	if err != nil {
		slog.Error("Failed to open", "error", err, "file path", o.filepath)
		return models.InventoryItem{}, myerrors.ErrFailOpenJson
	}

	var inventeryItems []models.InventoryItem
	if err := json.Unmarshal(byteValue, &inventeryItems); err != nil {
		slog.Error("Failed to unmarshal", "error", err)
		return models.InventoryItem{}, myerrors.ErrFailUnmarshal
	}

	for _, inventoryItem := range inventeryItems {
		if inventoryItem.IngredientID == id {
			return inventoryItem, nil
		}
	}

	return models.InventoryItem{}, myerrors.ErrNotFound
}

func (o *jsonInventoryRepository) CreateInventory(newInventoryItem models.InventoryItem) error {
	byteValue, err := utils.ReadFile(o.filepath)
	if err != nil {
		slog.Error("Failed to open", "error", err, "file path", o.filepath)
		return myerrors.ErrFailOpenJson
	}

	var inventeryItems []models.InventoryItem
	if err := json.Unmarshal(byteValue, &inventeryItems); err != nil {
		slog.Error("Failed to unmarshal", "error", err)
		return myerrors.ErrFailUnmarshal
	}
	inventeryItems = append(inventeryItems, newInventoryItem)

	filestring, _ := json.MarshalIndent(inventeryItems, "", "  ")
	os.WriteFile(*config.Dir+"/"+o.filepath, filestring, os.ModePerm)

	return nil
}

func (o *jsonInventoryRepository) UpdateInventory(id string, newInvItem models.InventoryItem) error {
	byteValue, err := utils.ReadFile(o.filepath)
	if err != nil {
		slog.Error("Failed to open", "error", err, "file path", o.filepath)
		return myerrors.ErrFailOpenJson

	}

	var inventeryItems []models.InventoryItem
	if err := json.Unmarshal(byteValue, &inventeryItems); err != nil {
		slog.Error("Failed to unmarshal", "error", err)
		return myerrors.ErrFailUnmarshal
	}

	var isFound bool

	for i := range inventeryItems {
		if inventeryItems[i].IngredientID == id {
			inventeryItems[i] = newInvItem
			isFound = true
		}
	}

	if !isFound {
		slog.Error("Failed to find", "error", myerrors.ErrNotFound)
		return myerrors.ErrNotFound
	}

	filestring, _ := json.MarshalIndent(inventeryItems, "", "  ")
	os.WriteFile(*config.Dir+"/"+o.filepath, filestring, os.ModePerm)

	return nil
}

func (o *jsonInventoryRepository) DeleteInventory(id string) error {
	byteValue, err := utils.ReadFile(o.filepath)
	if err != nil {
		slog.Error("Failed to open", "error", err, "file path", o.filepath)

		return myerrors.ErrFailOpenJson
	}

	var inventeryItems []models.InventoryItem
	if err := json.Unmarshal(byteValue, &inventeryItems); err != nil {
		slog.Error("Failed to unmarshal", "error", err)

		return myerrors.ErrFailUnmarshal
	}

	var isFound bool

	var newinventeryItems []models.InventoryItem
	for i := range inventeryItems {
		if inventeryItems[i].IngredientID == id {
			isFound = true
			continue
		}
		newinventeryItems = append(newinventeryItems, inventeryItems[i])
	}

	if !isFound {
		slog.Error("Failed to find", "error", myerrors.ErrNotFound)
		return myerrors.ErrNotFound
	}

	filestring, _ := json.MarshalIndent(newinventeryItems, "", "  ")
	os.WriteFile(*config.Dir+"/"+o.filepath, filestring, os.ModePerm)

	return nil
}
