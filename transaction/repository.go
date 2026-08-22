package transaction

type Repository interface {
	Add(Transaction) Transaction
	GetByID(userID int, id int) (Transaction, bool)
	GetFiltered(Filter) []Transaction
	Delete(userID int, id int) error
	Update(userID int, tx Transaction) (Transaction, error)
}
