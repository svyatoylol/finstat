package api

import (
	"encoding/json"
	"errors"
	"fin-stat/transaction"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Server struct {
	service *transaction.Service
}

type transactionRequest struct {
	Type     transaction.TransactionType `json:"type"`
	Amount   float64                     `json:"amount"`
	Category string                      `json:"category"`
}

type transactionResponse struct {
	ID       int64                       `json:"id"`
	UserID   int64                       `json:"user_id"`
	Type     transaction.TransactionType `json:"type"`
	Amount   float64                     `json:"amount"`
	Category string                      `json:"category"`
}

type balanceResponse struct {
	UserID  int64   `json:"user_id"`
	Balance float64 `json:"balance"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func NewHandler(service *transaction.Service) http.Handler {
	server := &Server{service: service}
	mux := http.NewServeMux()
	mux.HandleFunc("/transactions", server.transactions)
	mux.HandleFunc("/transactions/", server.transaction)
	mux.HandleFunc("/balance", server.balance)
	return mux
}

func (s *Server) transactions(w http.ResponseWriter, r *http.Request) {
	userID, err := userIDFromRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	switch r.Method {
	case http.MethodGet:
		s.listTransactions(w, r, userID)
	case http.MethodPost:
		s.createTransaction(w, r, userID)
	default:
		writeMethodNotAllowed(w, http.MethodGet, http.MethodPost)
	}
}

func (s *Server) transaction(w http.ResponseWriter, r *http.Request) {
	userID, err := userIDFromRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	id, err := strconv.ParseInt(strings.TrimPrefix(r.URL.Path, "/transactions/"), 10, 64)
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, errors.New("invalid transaction id"))
		return
	}

	switch r.Method {
	case http.MethodGet:
		t, err := s.service.GetByID(r.Context(), userID, id)
		if err != nil {
			s.writeServiceError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, toResponse(t))
	case http.MethodPut:
		s.updateTransaction(w, r, userID, id)
	case http.MethodDelete:
		if err := s.service.DeleteTransaction(r.Context(), userID, id); err != nil {
			s.writeServiceError(w, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	default:
		writeMethodNotAllowed(w, http.MethodGet, http.MethodPut, http.MethodDelete)
	}
}

func (s *Server) balance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeMethodNotAllowed(w, http.MethodGet)
		return
	}
	userID, err := userIDFromRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	balance, err := s.service.GetBalance(r.Context(), userID)
	if err != nil {
		s.writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, balanceResponse{UserID: userID, Balance: balance})
}

func (s *Server) listTransactions(w http.ResponseWriter, r *http.Request, userID int64) {
	limit, err := queryInt(r, "limit", transaction.DefaultPageSize)
	if err != nil || limit <= 0 {
		writeError(w, http.StatusBadRequest, errors.New("limit must be positive"))
		return
	}
	afterID, err := queryInt64(r, "after_id", 0)
	if err != nil || afterID < 0 {
		writeError(w, http.StatusBadRequest, errors.New("after_id must be non-negative"))
		return
	}

	filter := transaction.Filter{
		UserID:  userID,
		Limit:   limit,
		AfterID: afterID,
	}

	if typeValue := r.URL.Query().Get("type"); typeValue != "" {
		txType := transaction.TransactionType(typeValue)
		if txType != transaction.Incoming && txType != transaction.Outgoing {
			writeError(w, http.StatusBadRequest, transaction.ErrInvalidType)
			return
		}
		filter.Type = &txType
	}

	if category := r.URL.Query().Get("category"); category != "" {
		filter.Category = &category
	}

	transactions, err := s.service.GetFiltered(r.Context(), filter)
	if err != nil {
		s.writeServiceError(w, err)
		return
	}

	result := make([]transactionResponse, 0, len(transactions))
	for _, item := range transactions {
		result = append(result, toResponse(item))
	}
	writeJSON(w, http.StatusOK, result)
}

func (s *Server) createTransaction(w http.ResponseWriter, r *http.Request, userID int64) {
	var input transactionRequest
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	created, err := s.service.AddTransaction(r.Context(), transaction.Transaction{
		UserID: userID, Type: input.Type, Amount: input.Amount, Category: input.Category,
	})
	if err != nil {
		s.writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, toResponse(created))
}

func (s *Server) updateTransaction(w http.ResponseWriter, r *http.Request, userID, id int64) {
	var input transactionRequest
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	updated, err := s.service.UpdateTransaction(r.Context(), userID, transaction.Transaction{
		ID: id, UserID: userID, Type: input.Type, Amount: input.Amount, Category: input.Category,
	})
	if err != nil {
		s.writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toResponse(updated))
}

func userIDFromRequest(r *http.Request) (int64, error) {
	value := r.Header.Get("X-User-ID")
	if value == "" {
		return 0, errors.New("X-User-ID header is required")
	}
	id, err := strconv.ParseInt(value, 10, 64)
	if err != nil || id <= 0 {
		return 0, errors.New("X-User-ID must be a positive integer")
	}
	return id, nil
}

func queryInt(r *http.Request, name string, fallback int) (int, error) {
	value := r.URL.Query().Get(name)
	if value == "" {
		return fallback, nil
	}
	return strconv.Atoi(value)
}

func queryInt64(r *http.Request, name string, fallback int64) (int64, error) {
	value := r.URL.Query().Get(name)
	if value == "" {
		return fallback, nil
	}
	return strconv.ParseInt(value, 10, 64)
}

func decodeJSON(r *http.Request, value any) error {
	if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		return errors.New("Content-Type must be application/json")
	}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(value); err != nil {
		return errors.New("invalid JSON body")
	}
	return nil
}

func (s *Server) writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, transaction.ErrNotFound):
		writeError(w, http.StatusNotFound, err)
	case errors.Is(err, transaction.ErrInvalidAmount), errors.Is(err, transaction.ErrInvalidType):
		writeError(w, http.StatusBadRequest, err)
	default:
		writeError(w, http.StatusInternalServerError, errors.New("internal server error"))
	}
}

func toResponse(t transaction.Transaction) transactionResponse {
	return transactionResponse{ID: t.ID, UserID: t.UserID, Type: t.Type, Amount: t.Amount, Category: t.Category}
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, errorResponse{Error: err.Error()})
}

func writeMethodNotAllowed(w http.ResponseWriter, methods ...string) {
	w.Header().Set("Allow", strings.Join(methods, ", "))
	writeError(w, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
}
