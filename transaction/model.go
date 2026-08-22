package transaction

type TransactionType string

const (
	Incoming TransactionType = "incoming"
	Outgoing TransactionType = "outgoing"
)

type Transaction struct {
	ID       int
	UserID   int
	Type     TransactionType
	Amount   float64
	Category string
}
