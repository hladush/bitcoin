// Package handlers provides HTTP handlers for the Bitcoin tracker REST API
package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ihladush/bitcoin/internal/models"
	"github.com/ihladush/bitcoin/internal/services"
)

// BitcoinHandler handles HTTP requests for Bitcoin tracking
type BitcoinHandler struct {
	service *services.BitcoinService
}

// NewBitcoinHandler creates a new Bitcoin handler
func NewBitcoinHandler(service *services.BitcoinService) *BitcoinHandler {
	return &BitcoinHandler{service: service}
}

// AddAddress handles POST /addresses
func (h *BitcoinHandler) AddAddress(w http.ResponseWriter, r *http.Request) {
	var req models.AddAddressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Address == "" {
		h.writeError(w, http.StatusBadRequest, "Address is required")
		return
	}

	address, err := h.service.AddAddress(req.Address, req.Label)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.writeSuccess(w, http.StatusCreated, address)
}

// RemoveAddress handles DELETE /addresses/{address}
func (h *BitcoinHandler) RemoveAddress(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]

	if address == "" {
		h.writeError(w, http.StatusBadRequest, "Address parameter is required")
		return
	}

	if err := h.service.RemoveAddress(address); err != nil {
		h.writeError(w, http.StatusNotFound, err.Error())
		return
	}

	h.writeMessage(w, http.StatusOK, "Address removed successfully")
}

// GetAllAddresses handles GET /addresses
func (h *BitcoinHandler) GetAllAddresses(w http.ResponseWriter, r *http.Request) {
	addresses, err := h.service.GetAllAddresses()
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.writeSuccess(w, http.StatusOK, addresses)
}

// GetAddress handles GET /addresses/{address}
func (h *BitcoinHandler) GetAddress(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]

	if address == "" {
		h.writeError(w, http.StatusBadRequest, "Address parameter is required")
		return
	}

	addressWithBalance, err := h.service.GetAddress(address)
	if err != nil {
		h.writeError(w, http.StatusNotFound, err.Error())
		return
	}

	h.writeSuccess(w, http.StatusOK, addressWithBalance)
}

// GetBalance handles GET /addresses/{address}/balance
func (h *BitcoinHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]

	if address == "" {
		h.writeError(w, http.StatusBadRequest, "Address parameter is required")
		return
	}

	balance, err := h.service.GetBalance(address)
	if err != nil {
		h.writeError(w, http.StatusNotFound, err.Error())
		return
	}

	h.writeSuccess(w, http.StatusOK, balance)
}

// GetTransactions handles GET /addresses/{address}/transactions
func (h *BitcoinHandler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]

	if address == "" {
		h.writeError(w, http.StatusBadRequest, "Address parameter is required")
		return
	}

	// Parse pagination parameters
	limit := 50
	offset := 0

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	transactions, err := h.service.GetTransactions(address, limit, offset)
	if err != nil {
		h.writeError(w, http.StatusNotFound, err.Error())
		return
	}

	h.writeSuccess(w, http.StatusOK, transactions)
}

// SyncAddress handles POST /addresses/{address}/sync
func (h *BitcoinHandler) SyncAddress(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]

	if address == "" {
		h.writeError(w, http.StatusBadRequest, "Address parameter is required")
		return
	}

	if err := h.service.SyncAddress(address); err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.writeMessage(w, http.StatusOK, "Address synchronized successfully")
}

// SyncAllAddresses handles POST /sync
func (h *BitcoinHandler) SyncAllAddresses(w http.ResponseWriter, r *http.Request) {
	if err := h.service.SyncAllAddresses(); err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.writeMessage(w, http.StatusOK, "All addresses synchronized successfully")
}

// HealthCheck handles GET /health
func (h *BitcoinHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	h.writeSuccess(w, http.StatusOK, map[string]string{
		"status":  "healthy",
		"service": "bitcoin-tracker",
	})
}

// Helper methods for response handling
func (h *BitcoinHandler) writeSuccess(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(models.SuccessResponse(data))
}

func (h *BitcoinHandler) writeError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(models.ErrorResponse(message))
}

func (h *BitcoinHandler) writeMessage(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(models.MessageResponse(message))
}
