package transaction

type Filter struct {
	UserID   int
	Type     *TransactionType
	Category *string
}
