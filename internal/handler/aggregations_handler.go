package handler

import (
	"hot-coffee/internal/service"
	"hot-coffee/internal/utils/response"
	"net/http"

	myerrors "hot-coffee/internal/myErrors"
)

type AggregationsHandler interface {
	HandleGetSales(w http.ResponseWriter, r *http.Request)
	HandleGetPopItems(w http.ResponseWriter, r *http.Request)
}

type aggregationsHandler struct {
	service service.AggregationsService
}

func NewAggregationsHandler(service service.AggregationsService) AggregationsHandler {
	return &aggregationsHandler{service: service}
}

// Get the total sales amount.
func (s *aggregationsHandler) HandleGetSales(w http.ResponseWriter, r *http.Request) {
	byteValue, err := s.service.ServiceGetTotal()
	if err != nil {
		response.SendError(w, http.StatusInternalServerError, "Failed to retrieve total sales amount", nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(byteValue)
}

// Get a list of popular menu items.
func (s *aggregationsHandler) HandleGetPopItems(w http.ResponseWriter, r *http.Request) {
	byteValue, err := s.service.ServiceGetPopular()
	if err == myerrors.ErrNoItems {
		response.SendError(w, http.StatusNotFound, "Failed to retrieve popular items", err)
	} else if err != nil {
		response.SendError(w, http.StatusInternalServerError, "Failed to retrieve popular items", nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(byteValue)
}
