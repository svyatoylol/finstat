package transaction

type memoryRepository struct {
	data []Transaction
}

func NewMemoryRepository() *memoryRepository {
	return &memoryRepository{
		data: []Transaction{},
	}
}

func (r *memoryRepository) Add(t Transaction) {
	r.data = append(r.data, t)
}

func (r *memoryRepository) GetAll() []Transaction {
	return r.data
}
