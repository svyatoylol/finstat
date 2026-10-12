package transaction

type TransactionType string

const (
	Incoming TransactionType = "incoming"
	Outgoing TransactionType = "outgoing"
)

type Transaction struct {
	ID       int64
	UserID   int64
	Type     TransactionType
	Amount   float64
	Category string
}
