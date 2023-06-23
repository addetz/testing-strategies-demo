package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/addetz/testing-strategies-demo/data"
	"github.com/gorilla/mux"
)

type ResponseType interface {
	data.Events | data.Talks | ErrorResponse
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type Handler struct {
	eventService *data.EventService
}

func NewHandler(es *data.EventService) *Handler {
	return &Handler{
		eventService: es,
	}
}

func (h *Handler) GetEventsHandler(w http.ResponseWriter, r *http.Request) {
	events := h.eventService.GetEvents()
	writeResponse[data.Events](w, http.StatusOK, &events)
}

func (h *Handler) GetEventTalksHandler(w http.ResponseWriter, r *http.Request) {
	eventID := mux.Vars(r)["id"]
	day := r.URL.Query().Get("day")
	var talks *data.Talks
	if len(day) != 0 {
		h.fetchFilteredEvents(w, eventID, day)
		return
	}
	talks, err := h.eventService.GetEventTalks(eventID)
	if err != nil {
		writeResponse[ErrorResponse](w, http.StatusBadRequest, &ErrorResponse{
			Error: fmt.Errorf("GetEventTalksHandler:%v", err).Error(),
		})
		return
	}
	writeResponse[data.Talks](w, http.StatusOK, talks)
}

func (h *Handler) fetchFilteredEvents(w http.ResponseWriter, eventID, day string) {
	parsedDay, err := strconv.Atoi(day)
	if err != nil {
		writeResponse[ErrorResponse](w, http.StatusBadRequest, &ErrorResponse{
			Error: fmt.Errorf("GetEventTalksHandler:%v", err).Error(),
		})
		return
	}
	talks, err := h.eventService.GetEventFilteredTalks(eventID, parsedDay)
	if err != nil {
		writeResponse[ErrorResponse](w, http.StatusBadRequest, &ErrorResponse{
			Error: fmt.Errorf("GetEventTalksHandler:%v", err).Error(),
		})
		return
	}
	writeResponse[data.Talks](w, http.StatusOK, talks)
}

// writeResponse is a helper method that allows to write the HTTP status & response
func writeResponse[T ResponseType](w http.ResponseWriter, status int, resp *T) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if status != http.StatusOK {
		w.WriteHeader(status)
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error encoding resp %v:%s", resp, err)
	}
}
