package transaction

import "fmt"

type memoryRepository struct {
	data   []Transaction
	nextID int
}

func NewMemoryRepository() *memoryRepository {
	return &memoryRepository{
		data:   []Transaction{},
		nextID: 1,
	}
}

func (r *memoryRepository) Add(t Transaction) Transaction {
	t.ID = r.nextID
	r.nextID++
	r.data = append(r.data, t)
	return t
}

func (r *memoryRepository) GetByID(userID int, id int) (Transaction, bool) {
	for _, t := range r.data {
		if t.ID == id && t.UserID == userID {
			return t, true
		}
	}
	return Transaction{}, false
}

func (r *memoryRepository) GetFiltered(f Filter) []Transaction {
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
	return filtered
}

func (r *memoryRepository) Delete(userid int, id int) error {
	for i, t := range r.data {
		if t.ID == id && t.UserID == userid {
			r.data = append(r.data[:i], r.data[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("transaction with id %d not found", id)
}

func (r *memoryRepository) Update(userID int, t Transaction) (Transaction, error) {
	for i, tx := range r.data {
		if tx.ID == t.ID && tx.UserID == userID {
			t.UserID = userID
			r.data[i] = t
			return t, nil
		}
	}
	return Transaction{}, fmt.Errorf("transaction with id %d not found", t.ID)
}
