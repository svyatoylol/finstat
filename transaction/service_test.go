package transaction

import (
	"context"
	"testing"
)

func TestServiceGetBalance(t *testing.T) {
	repo := NewMemoryRepository()
	service := NewService(repo)
	ctx := context.Background()

	for _, tx := range []Transaction{
		{UserID: 7, Type: Incoming, Amount: 250.50, Category: "salary"},
		{UserID: 7, Type: Outgoing, Amount: 50.25, Category: "food"},
		{UserID: 8, Type: Incoming, Amount: 1000, Category: "salary"},
	} {
		if _, err := service.AddTransaction(ctx, tx); err != nil {
			t.Fatalf("add transaction: %v", err)
		}
	}

	balance, err := service.GetBalance(ctx, 7)
	if err != nil {
		t.Fatalf("get balance: %v", err)
	}
	if balance != 200.25 {
		t.Fatalf("expected balance 200.25, got %.2f", balance)
	}
}

func TestMemoryRepositoryGetFilteredPagination(t *testing.T) {
	repo := NewMemoryRepository()
	ctx := context.Background()

	for i := 0; i < 4; i++ {
		if _, err := repo.Add(ctx, Transaction{
			UserID: 1,
			Type:   Incoming,
			Amount: float64(i + 1),
		}); err != nil {
			t.Fatalf("add transaction %d: %v", i, err)
		}
	}

	page, err := repo.GetFiltered(ctx, Filter{
		UserID: 1,
		Type:   pointer(Incoming),
		Limit:  2,
	})
	if err != nil {
		t.Fatalf("get first page: %v", err)
	}
	if len(page) != 2 || page[0].ID != 1 || page[1].ID != 2 {
		t.Fatalf("unexpected first page: %#v", page)
	}

	page, err = repo.GetFiltered(ctx, Filter{
		UserID:  1,
		Type:    pointer(Incoming),
		Limit:   2,
		AfterID: 2,
	})
	if err != nil {
		t.Fatalf("get second page: %v", err)
	}
	if len(page) != 2 || page[0].ID != 3 || page[1].ID != 4 {
		t.Fatalf("unexpected second page: %#v", page)
	}
}

func pointer[T any](value T) *T {
	return &value
}
