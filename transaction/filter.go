package transaction

const DefaultPageSize = 100

type Filter struct {
	UserID   int64
	Type     *TransactionType
	Category *string
	Limit    int
	AfterID  int64
}
