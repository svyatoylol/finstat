package api

import (
	"bytes"
	"encoding/json"
	"fin-stat/transaction"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFullUserScenario(t *testing.T) {
	service := transaction.NewService(transaction.NewMemoryRepository())
	handler := NewHandler(service)

	doReq := func(method, path, body, userID string) *httptest.ResponseRecorder {
		var bodyReader io.Reader
		if body != "" {
			bodyReader = bytes.NewBufferString(body)
		}
		req := httptest.NewRequest(method, path, bodyReader)
		if userID != "" {
			req.Header.Set("X-User-ID", userID)
		}
		if body != "" {
			req.Header.Set("Content-Type", "application/json")
		}
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		return rec
	}

	// 1. Запрос без X-User-ID должен возвращать 400
	resp := doReq(http.MethodGet, "/transactions", "", "")
	if resp.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.Code)
	}

	// 2. Добавление дохода (+50000)
	resp = doReq(http.MethodPost, "/transactions", `{"type":"incoming","amount":50000,"category":"salary"}`, "1")
	if resp.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.Code)
	}
	var t1 transactionResponse
	_ = json.NewDecoder(resp.Body).Decode(&t1)
	if t1.ID != 1 || t1.Amount != 50000 {
		t.Fatalf("unexpected transaction: %+v", t1)
	}

	// 3. Добавление расхода (-12000)
	resp = doReq(http.MethodPost, "/transactions", `{"type":"outgoing","amount":12000,"category":"rent"}`, "1")
	if resp.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.Code)
	}

	// 4. Проверка баланса (50000 - 12000 = 38000)
	resp = doReq(http.MethodGet, "/balance", "", "1")
	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.Code)
	}
	var bal balanceResponse
	_ = json.NewDecoder(resp.Body).Decode(&bal)
	if bal.Balance != 38000 {
		t.Fatalf("expected balance 38000, got %.2f", bal.Balance)
	}

	// 5. Валидация: отрицательная сумма
	resp = doReq(http.MethodPost, "/transactions", `{"type":"incoming","amount":-100,"category":"test"}`, "1")
	if resp.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for negative amount, got %d", resp.Code)
	}

	// 6. Удаление транзакции #2 (расход)
	resp = doReq(http.MethodDelete, "/transactions/2", "", "1")
	if resp.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", resp.Code)
	}

	// 7. Баланс после удаления расхода должен снова стать 50000
	resp = doReq(http.MethodGet, "/balance", "", "1")
	var balAfter balanceResponse
	_ = json.NewDecoder(resp.Body).Decode(&balAfter)
	if balAfter.Balance != 50000 {
		t.Fatalf("expected balance 50000, got %.2f", balAfter.Balance)
	}

	// 8. Несуществующая транзакция -> 404
	resp = doReq(http.MethodGet, "/transactions/999", "", "1")
	if resp.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", resp.Code)
	}

	// 9. Неподдерживаемый метод -> 405
	resp = doReq(http.MethodPatch, "/transactions/1", `{}`, "1")
	if resp.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", resp.Code)
	}
}
