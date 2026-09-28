package transaction

import (
	"context"
	"errors"
)

var ErrNotFound = errors.New("transaction not found")

type Repository interface {
	Add(ctx context.Context, t Transaction) (Transaction, error)
	GetByID(ctx context.Context, userID, id int) (Transaction, error)
	GetFiltered(ctx context.Context, f Filter) ([]Transaction, error)
	Delete(ctx context.Context, userID, id int) error
	Update(ctx context.Context, userID int, t Transaction) (Transaction, error)
}
