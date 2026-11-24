package api

import (
	"bytes"
	"encoding/json"
	"fin-stat/transaction"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTransactionsCRUD(t *testing.T) {
	service := transaction.NewService(transaction.NewMemoryRepository())
	handler := NewHandler(service)

	create := request(t, handler, http.MethodPost, "/transactions", `{"type":"incoming","amount":100,"category":"salary"}`)
	if create.Code != http.StatusCreated {
		t.Fatalf("create status: got %d, want %d", create.Code, http.StatusCreated)
	}

	var created transactionResponse
	decode(t, create, &created)
	if created.ID != 1 || created.UserID != 7 {
		t.Fatalf("unexpected created transaction: %+v", created)
	}

	list := request(t, handler, http.MethodGet, "/transactions?limit=10", "")
	if list.Code != http.StatusOK {
		t.Fatalf("list status: got %d, want %d", list.Code, http.StatusOK)
	}

	update := request(t, handler, http.MethodPut, "/transactions/1", `{"type":"incoming","amount":125,"category":"bonus"}`)
	if update.Code != http.StatusOK {
		t.Fatalf("update status: got %d, want %d", update.Code, http.StatusOK)
	}

	remove := request(t, handler, http.MethodDelete, "/transactions/1", "")
	if remove.Code != http.StatusNoContent {
		t.Fatalf("delete status: got %d, want %d", remove.Code, http.StatusNoContent)
	}
}

func TestTransactionsRequireUserID(t *testing.T) {
	handler := NewHandler(transaction.NewService(transaction.NewMemoryRepository()))
	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/transactions", nil)
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("status: got %d, want %d", response.Code, http.StatusBadRequest)
	}
}

func TestListTransactionsFiltering(t *testing.T) {
	service := transaction.NewService(transaction.NewMemoryRepository())
	handler := NewHandler(service)

	request(t, handler, http.MethodPost, "/transactions", `{"type":"incoming","amount":1000,"category":"salary"}`)
	request(t, handler, http.MethodPost, "/transactions", `{"type":"incoming","amount":200,"category":"bonus"}`)
	request(t, handler, http.MethodPost, "/transactions", `{"type":"outgoing","amount":50,"category":"food"}`)
	request(t, handler, http.MethodPost, "/transactions", `{"type":"outgoing","amount":300,"category":"rent"}`)

	// Filter by both type AND category
	both := request(t, handler, http.MethodGet, "/transactions?type=incoming&category=salary", "")
	if both.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d", both.Code, http.StatusOK)
	}
	var bothList []transactionResponse
	decode(t, both, &bothList)
	if len(bothList) != 1 || bothList[0].Category != "salary" || bothList[0].Type != transaction.Incoming {
		t.Fatalf("expected 1 incoming salary, got %+v", bothList)
	}

	// Filter by type only
	typeOnly := request(t, handler, http.MethodGet, "/transactions?type=outgoing", "")
	var typeList []transactionResponse
	decode(t, typeOnly, &typeList)
	if len(typeList) != 2 {
		t.Fatalf("expected 2 outgoing transactions, got %d", len(typeList))
	}

	// Filter by category only
	catOnly := request(t, handler, http.MethodGet, "/transactions?category=bonus", "")
	var catList []transactionResponse
	decode(t, catOnly, &catList)
	if len(catList) != 1 || catList[0].Amount != 200 {
		t.Fatalf("expected 1 bonus transaction, got %+v", catList)
	}

	// Invalid type
	invalid := request(t, handler, http.MethodGet, "/transactions?type=unknown", "")
	if invalid.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid type, got %d", invalid.Code)
	}
}

func request(t *testing.T, handler http.Handler, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, path, bytes.NewBufferString(body))
	req.Header.Set("X-User-ID", "7")
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, req)
	return response
}

func decode(t *testing.T, response *httptest.ResponseRecorder, value any) {
	t.Helper()
	if err := json.NewDecoder(response.Body).Decode(value); err != nil {
		t.Fatalf("decode response: %v", err)
	}
}
