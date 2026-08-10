package transaction

type Repository interface {
	Add(Transaction) Transaction
	GetAll() []Transaction
}
