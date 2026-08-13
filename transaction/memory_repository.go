package transaction

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

func (r *memoryRepository) GetAll() []Transaction {
	return r.data
}

func (r *memoryRepository) GetByID(id int) (Transaction, bool) {
	for _, t := range r.data {
		if t.ID == id {
			return t, true
		}
	}
	return Transaction{}, false
}

func (r *memoryRepository) GetFiltered(f Filter) []Transaction {
	var filtered []Transaction
	for _, tx := range r.data {
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
