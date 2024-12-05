package handler

import (
	"io"
	"net/http"

	myerrors "hot-coffee/internal/myErrors"
	"hot-coffee/internal/service"
	"hot-coffee/internal/utils/response"
	"hot-coffee/internal/utils/validation"
)

type InventoryHandler interface {
	HandleGetInventory(w http.ResponseWriter, r *http.Request)
	HandleGetInventoryID(w http.ResponseWriter, r *http.Request)
	HandlePostInventory(w http.ResponseWriter, r *http.Request)
	HandlePutInventoryID(w http.ResponseWriter, r *http.Request)
	HandleDeleteInventory(w http.ResponseWriter, r *http.Request)
}

type inventoryHandler struct {
	service service.InventoryService
}

func NewInventoryHandler(service service.InventoryService) InventoryHandler {
	return &inventoryHandler{service: service}
}

// Retrieve all inventory items.
func (s *inventoryHandler) HandleGetInventory(w http.ResponseWriter, r *http.Request) {
	byteValue, err := s.service.ServiceGetInventory()
	if err != nil {
		response.SendError(w, http.StatusInternalServerError, "Failed to retrieve inventory", nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(byteValue)
}

// Retrieve a specific inventory item.
func (s *inventoryHandler) HandleGetInventoryID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	byteValue, err := s.service.ServiceGetInventoryID(id)

	if err == myerrors.ErrNotFound {
		response.SendError(w, http.StatusNotFound, "Failed to retrieve inventory", err)
		return
	} else if err != nil {
		response.SendError(w, http.StatusInternalServerError, "Failed to retrieve inventory", nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(byteValue)
}

// Add a new inventory item.
func (s *inventoryHandler) HandlePostInventory(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if !validation.IsJSON(contentType) {
		response.SendError(w, http.StatusBadRequest, "Not a JSON", nil)
		return
	}

	inventoryByte, err := io.ReadAll(r.Body)
	if err != nil {
		response.SendError(w, http.StatusInternalServerError, "Failed to create inventory", nil)
		return
	}

	err = s.service.ServiceCreateInventory(inventoryByte)
	switch err {
	case myerrors.ErrIdRequired,
		myerrors.ErrNameRequired,
		myerrors.ErrInvalidQuantity,
		myerrors.ErrUnitRequired:
		response.SendError(w, http.StatusBadRequest, "Failed to create inventory", err)
		return
	default:
		if err != nil {
			response.SendError(w, http.StatusInternalServerError, "Failed to create inventory", myerrors.ErrInvalidJson)
			return
		}
	}

	response.SendMessage(w, http.StatusOK, "inventory item succesfuly added")
}

// Update an inventory item.
func (s *inventoryHandler) HandlePutInventoryID(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if !validation.IsJSON(contentType) {
		response.SendError(w, http.StatusBadRequest, "Not a JSON", nil)
		return
	}

	id := r.PathValue("id")

	inventoryByte, err := io.ReadAll(r.Body)
	if err != nil {
		response.SendError(w, http.StatusInternalServerError, "Failed to update inventory", nil)
	}

	err = s.service.ServiceUpdateInventory(id, inventoryByte)
	switch err {
	case myerrors.ErrIdRequired,
		myerrors.ErrNameRequired,
		myerrors.ErrInvalidQuantity,
		myerrors.ErrUnitRequired:
		response.SendError(w, http.StatusBadRequest, "Failed to update inventory", err)
		return
	case myerrors.ErrNotFound:
		response.SendError(w, http.StatusNotFound, "Failed to update inventory", err)
		return
	default:
		if err != nil {
			response.SendError(w, http.StatusInternalServerError, "Failed to update inventory", myerrors.ErrInvalidJson)
			return
		}
	}

	response.SendMessage(w, http.StatusCreated, "inventory item succesfuly updated")
}

// Delete an inventory item.
func (s *inventoryHandler) HandleDeleteInventory(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := s.service.ServiceDeleteInventory(id)
	switch err {
	case myerrors.ErrNotFound:
		response.SendError(w, http.StatusNotFound, "Failed to delete inventory", err)
		return
	default:
		if err != nil {
			response.SendError(w, http.StatusInternalServerError, "Failed to delete inventory", nil)
			return
		}
	}

	response.SendMessage(w, http.StatusAccepted, "inventory item succesfuly deleted")
}
