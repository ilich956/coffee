package handler

import (
	"hot-coffee/internal/service"
	"hot-coffee/internal/utils/response"
	"hot-coffee/internal/utils/validation"
	"io"
	"net/http"

	myerrors "hot-coffee/internal/myErrors"
)

type OrderHandler interface {
	HandleGetOrder(w http.ResponseWriter, r *http.Request)
	HandleGetOrderID(w http.ResponseWriter, r *http.Request)
	HandlePostOrder(w http.ResponseWriter, r *http.Request)
	HandlePostOrderClose(w http.ResponseWriter, r *http.Request)
	HandlePutOrderID(w http.ResponseWriter, r *http.Request)
	HandleDeleteOrder(w http.ResponseWriter, r *http.Request)
}

type orderHandler struct {
	service service.OrderService
}

func NewOrderHandler(service service.OrderService) OrderHandler {
	return &orderHandler{service: service}
}

// Retrieve all orders.
func (s *orderHandler) HandleGetOrder(w http.ResponseWriter, r *http.Request) {
	byteValue, err := s.service.ServiceGetOrder()
	if err != nil {
		response.SendError(w, http.StatusInternalServerError, "Failed to retrieve orders", nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(byteValue)
}

// Retrieve a specific order by ID.
func (s *orderHandler) HandleGetOrderID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	byteValue, err := s.service.ServiceGetOrderID(id)

	if err == myerrors.ErrNotFound {
		response.SendError(w, http.StatusNotFound, "Failed to retrieve order", err)
		return
	} else if err != nil {
		response.SendError(w, http.StatusInternalServerError, "Failed to retrieve order", nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(byteValue)
}

// Create a new order.
func (s *orderHandler) HandlePostOrder(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if !validation.IsJSON(contentType) {
		response.SendError(w, http.StatusBadRequest, "Not a JSON", nil)
		return
	}

	orderByte, err := io.ReadAll(r.Body)
	if err != nil {
		response.SendError(w, http.StatusInternalServerError, "Failed to create order", nil)
		return
	}

	err = s.service.ServicePostOrder(orderByte)
	switch err {
	case myerrors.ErrNameRequired,
		myerrors.ErrItemsRequired,
		myerrors.ErrIdRequired,
		myerrors.ErrEmptyOrder,
		myerrors.ErrAbsentItem,
		myerrors.ErrInvalidQuantity:
		response.SendError(w, http.StatusBadRequest, "Failed to create order", err)
		return
		////////////////////////////////////////////////////////////////////////////////////////////
	default:
		if err != nil {
			response.SendError(w, http.StatusInternalServerError, "Failed to create order", myerrors.ErrInvalidJson)
			return
		}
	}

	response.SendMessage(w, http.StatusOK, "order succesfuly created")
}

// Close an order.
func (s *orderHandler) HandlePostOrderClose(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := s.service.ServicePostOrderClose(id)

	switch err {
	case myerrors.ErrOrderClosed:
		response.SendError(w, http.StatusConflict, "Failed to close order", err)
		return
	case myerrors.ErrNotFound:
		response.SendError(w, http.StatusNotFound, "Failed to close order", err)
		return
	default:
		if err != nil {
			response.SendError(w, http.StatusInternalServerError, "Failed to close order", err)
			return
		}
	}
	response.SendMessage(w, http.StatusCreated, "order succesfuly closed")
}

// Update an existing order.
func (s *orderHandler) HandlePutOrderID(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if !validation.IsJSON(contentType) {
		response.SendError(w, http.StatusBadRequest, "Not a JSON", nil)
		return
	}

	id := r.PathValue("id")

	orderByte, err := io.ReadAll(r.Body)
	if err != nil {
		response.SendError(w, http.StatusInternalServerError, "Failed to update an order", nil)
	}

	err = s.service.ServicePutOrderID(id, orderByte)

	switch err {
	case myerrors.ErrNameRequired,
		myerrors.ErrItemsRequired,
		myerrors.ErrIdRequired,
		myerrors.ErrInvalidQuantity,
		myerrors.ErrOrderClosed:
		response.SendError(w, http.StatusBadRequest, "Failed to update an order", err)
		return
	case myerrors.ErrNotFound:
		response.SendError(w, http.StatusNotFound, "Failed to update an order", err)
		return
	default:
		if err != nil {
			response.SendError(w, http.StatusInternalServerError, "Failed to update an order", myerrors.ErrInvalidJson)
			return
		}
	}

	response.SendMessage(w, http.StatusCreated, "order succesfuly updated")
}

// Delete an order.
func (s *orderHandler) HandleDeleteOrder(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	err := s.service.ServiceDeleteOrder(id)
	switch err {
	case myerrors.ErrNotFound:
		response.SendError(w, http.StatusNotFound, "Failed to delete an order", err)
		return
	default:
		if err != nil {
			response.SendError(w, http.StatusInternalServerError, "Failed to delete an order", nil)
			return
		}
	}
	response.SendMessage(w, http.StatusAccepted, "order succesfuly deleted")
}
