package service

import (
	"encoding/json"
	"hot-coffee/internal/dal"
	"hot-coffee/internal/utils/validation"
	"hot-coffee/models"
	"log/slog"

	myerrors "hot-coffee/internal/myErrors"
)

type InventoryService interface {
	ServiceGetInventory() ([]byte, error)
	ServiceGetInventoryID(id string) ([]byte, error)
	ServiceCreateInventory(newInventoryItem []byte) error
	ServiceUpdateInventory(id string, newInventoryItem []byte) error
	ServiceDeleteInventory(id string) error
}

type inventoryService struct {
	repo dal.InventoryRepository
}

func NewInventoryService(repo dal.InventoryRepository) InventoryService {
	return &inventoryService{repo: repo}
}

func (i *inventoryService) ServiceGetInventory() ([]byte, error) {
	inventoryStruct, err := i.repo.GetInventory()
	if err != nil {
		return nil, err
	}

	jsonFile, err := json.MarshalIndent(inventoryStruct, "", "  ")
	if err != nil {
		slog.Error("Failed to marshal", "error", err)
		return nil, myerrors.ErrFailMarshal
	}
	return jsonFile, nil
}

func (i *inventoryService) ServiceGetInventoryID(id string) ([]byte, error) {
	inventoryStruct, err := i.repo.GetInventoryID(id)
	if err == myerrors.ErrNotFound {
		slog.Error("Failed to find", "error", myerrors.ErrNotFound)
	}
	if err != nil {
		return nil, err
	}

	jsonFile, err := json.MarshalIndent(inventoryStruct, "", "  ")
	if err != nil {
		slog.Error("Failed to marshal", "error", err)
		return nil, myerrors.ErrFailMarshal
	}
	return jsonFile, nil
}

func (i *inventoryService) ServiceCreateInventory(newInventoryItem []byte) error {
	var inventory models.InventoryItem
	if err := json.Unmarshal(newInventoryItem, &inventory); err != nil {
		slog.Error("Failed to unmarshal", "error", err)
		return myerrors.ErrFailUnmarshal
	}

	if err := validation.CheckInventory(inventory); err != nil {
		return err
	}

	checkInventoryID, _ := i.repo.GetInventoryID(inventory.IngredientID)
	if checkInventoryID.IngredientID == inventory.IngredientID {
		slog.Error("Failed to create inventory item", "error", myerrors.ErrIDExist)
		return myerrors.ErrIDExist
	}

	return i.repo.CreateInventory(inventory)
}

func (i *inventoryService) ServiceUpdateInventory(id string, newInventoryItem []byte) error {
	var inventory models.InventoryItem
	if err := json.Unmarshal(newInventoryItem, &inventory); err != nil {
		slog.Error("Failed to unmarshal", "error", err)
		return myerrors.ErrFailUnmarshal
	}

	if err := validation.CheckInventory(inventory); err != nil {
		return err
	}

	return i.repo.UpdateInventory(id, inventory)
}

func (i *inventoryService) ServiceDeleteInventory(id string) error {
	return i.repo.DeleteInventory(id)
}
