package transaction

import (
	"context"
	"sync"
)

var _ Repository = (*memoryRepository)(nil)

type memoryRepository struct {
	mu     sync.Mutex
	data   []Transaction
	nextID int
}

func NewMemoryRepository() *memoryRepository {
	return &memoryRepository{nextID: 1}
}

func (r *memoryRepository) Add(_ context.Context, t Transaction) (Transaction, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	t.ID = r.nextID
	r.nextID++
	r.data = append(r.data, t)
	return t, nil
}

func (r *memoryRepository) GetByID(_ context.Context, userID, id int) (Transaction, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, t := range r.data {
		if t.ID == id && t.UserID == userID {
			return t, nil
		}
	}
	return Transaction{}, ErrNotFound
}

func (r *memoryRepository) GetFiltered(_ context.Context, f Filter) ([]Transaction, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var filtered []Transaction
	for _, tx := range r.data {
		if tx.UserID != f.UserID {
			continue
		}
		if f.Type != nil && tx.Type != *f.Type {
			continue
		}
		if f.Category != nil && tx.Category != *f.Category {
			continue
		}
		filtered = append(filtered, tx)
	}
	return filtered, nil
}

func (r *memoryRepository) Delete(_ context.Context, userID, id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, t := range r.data {
		if t.ID == id && t.UserID == userID {
			r.data = append(r.data[:i], r.data[i+1:]...)
			return nil
		}
	}
	return ErrNotFound
}

func (r *memoryRepository) Update(_ context.Context, userID int, t Transaction) (Transaction, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, tx := range r.data {
		if tx.ID == t.ID && tx.UserID == userID {
			t.UserID = userID
			r.data[i] = t
			return t, nil
		}
	}
	return Transaction{}, ErrNotFound
}
