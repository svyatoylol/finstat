package transaction

type Repository interface {
	Add(Transaction) Transaction
	GetByID(int) (Transaction, bool)
	GetFiltered(Filter) []Transaction
	Delete(int) error
	Update(Transaction) (Transaction, error)
}
