package transaction

type Repository interface {
	Add(Transaction) Transaction
	GetAll() []Transaction
	GetByID(int) (Transaction, bool)
	GetFiltered(Filter) []Transaction
}
