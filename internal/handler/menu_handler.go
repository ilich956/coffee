package handler

import (
	"io"
	"net/http"

	myerrors "hot-coffee/internal/myErrors"
	"hot-coffee/internal/service"
	"hot-coffee/internal/utils/response"
	"hot-coffee/internal/utils/validation"
)

type MenuHandler interface {
	HandleGetMenu(w http.ResponseWriter, r *http.Request)
	HandleGetMenuID(w http.ResponseWriter, r *http.Request)
	HandlePostMenu(w http.ResponseWriter, r *http.Request)
	HandlePutMenuID(w http.ResponseWriter, r *http.Request)
	HandleDeleteMenuID(w http.ResponseWriter, r *http.Request)
}

type menuHandler struct {
	service service.MenuService
}

func NewMenuHandler(service service.MenuService) MenuHandler {
	return &menuHandler{service: service}
}

// Retrieve all menu items.
func (s *menuHandler) HandleGetMenu(w http.ResponseWriter, r *http.Request) {
	byteValue, err := s.service.ServiceGetMenu()
	if err != nil {
		response.SendError(w, http.StatusInternalServerError, "Failed to retrieve menu", nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(byteValue)
}

// Retrieve a specific menu item.
func (s *menuHandler) HandleGetMenuID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	byteValue, err := s.service.ServiceGetMenuID(id)

	if err == myerrors.ErrNotFound {
		response.SendError(w, http.StatusNotFound, "Failed to retrieve menu", err)
		return
	} else if err != nil {
		response.SendError(w, http.StatusInternalServerError, "Failed to retrieve menu", nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(byteValue)
}

// Add a new menu item.
func (s *menuHandler) HandlePostMenu(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if !validation.IsJSON(contentType) {
		response.SendError(w, http.StatusBadRequest, "Not a JSON", nil)
		return
	}

	menuByte, err := io.ReadAll(r.Body)
	if err != nil {
		response.SendError(w, http.StatusInternalServerError, "Failed to create menu", nil)
		return
	}

	err = s.service.ServiceCreateMenu(menuByte)
	switch err {
	case myerrors.ErrNameRequired,
		myerrors.ErrIdRequired,
		myerrors.ErrInvalidQuantity,
		myerrors.ErrDescriptionRequired,
		myerrors.ErrPriceRequired,
		myerrors.ErrIngredientsRequired:
		response.SendError(w, http.StatusBadRequest, "Failed to create menu", err)
		return
		////////////////////////////////////////////////////////////////////////////////////////////
	default:
		if err != nil {
			response.SendError(w, http.StatusInternalServerError, "Failed to create menu", myerrors.ErrInvalidJson)
			return
		}
	}

	response.SendMessage(w, http.StatusOK, "menu succesfuly created")
}

// Update a menu item.
func (s *menuHandler) HandlePutMenuID(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if !validation.IsJSON(contentType) {
		response.SendError(w, http.StatusBadRequest, "Not a JSON", nil)
		return
	}

	id := r.PathValue("id")

	menuByte, err := io.ReadAll(r.Body)
	if err != nil {
		response.SendError(w, http.StatusInternalServerError, "Failed to update menu", nil)
	}

	err = s.service.ServiceUpdateMenu(id, menuByte)

	switch err {
	case myerrors.ErrNameRequired,
		myerrors.ErrIdRequired,
		myerrors.ErrInvalidQuantity,
		myerrors.ErrDescriptionRequired,
		myerrors.ErrPriceRequired,
		myerrors.ErrIngredientsRequired:
		response.SendError(w, http.StatusBadRequest, "Failed to update an menu", err)
		return
	case myerrors.ErrNotFound:
		response.SendError(w, http.StatusNotFound, "Failed to update an menu", err)
		return
	default:
		if err != nil {
			response.SendError(w, http.StatusInternalServerError, "Failed to update an menu", myerrors.ErrInvalidJson)
			return
		}
	}
	response.SendMessage(w, http.StatusCreated, "menu succesfuly updated")
}

// Delete a menu item.
func (s *menuHandler) HandleDeleteMenuID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := s.service.ServiceDeleteMenu(id)
	switch err {
	case myerrors.ErrNotFound:
		response.SendError(w, http.StatusNotFound, "Failed to delete menu", err)
		return
	default:
		if err != nil {
			response.SendError(w, http.StatusInternalServerError, "Failed to delete menu", nil)
			return
		}
	}

	response.SendMessage(w, http.StatusAccepted, "menu succesfuly deleted")
}
