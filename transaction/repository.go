package transaction

import (
	"context"
	"errors"
)

var ErrNotFound = errors.New("transaction not found")

type Repository interface {
	Add(ctx context.Context, t Transaction) (Transaction, error)
	GetByID(ctx context.Context, userID, id int64) (Transaction, error)
	GetFiltered(ctx context.Context, f Filter) ([]Transaction, error)
	GetBalance(ctx context.Context, userID int64) (float64, error)
	Delete(ctx context.Context, userID, id int64) error
	Update(ctx context.Context, userID int64, t Transaction) (Transaction, error)
}
