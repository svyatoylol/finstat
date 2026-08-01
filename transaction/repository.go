package transaction

type Repository interface {
	Add(Transaction)
	GetAll() []Transaction
}
