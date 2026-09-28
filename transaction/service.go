package transaction

import (
	"context"
	"errors"
	"fmt"
)

var (
	ErrInvalidAmount = errors.New("amount must be positive")
	ErrInvalidType   = errors.New("invalid transaction type")
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func validate(t Transaction) error {
	if t.Amount <= 0 {
		return ErrInvalidAmount
	}
	if t.Type != Incoming && t.Type != Outgoing {
		return ErrInvalidType
	}
	return nil
}

func (s *Service) GetBalance(ctx context.Context, userID int) (float64, error) {
	txs, err := s.repo.GetFiltered(ctx, Filter{UserID: userID})
	if err != nil {
		return 0, fmt.Errorf("get balance: %w", err)
	}

	var balance float64
	for _, t := range txs {
		switch t.Type {
		case Incoming:
			balance += t.Amount
		case Outgoing:
			balance -= t.Amount
		}
	}
	return balance, nil
}

func (s *Service) AddTransaction(ctx context.Context, t Transaction) (Transaction, error) {
	if err := validate(t); err != nil {
		return Transaction{}, err
	}
	return s.repo.Add(ctx, t)
}

func (s *Service) GetAllByType(ctx context.Context, userID int, tType TransactionType) ([]Transaction, error) {
	return s.repo.GetFiltered(ctx, Filter{
		UserID: userID,
		Type:   &tType,
	})
}

func (s *Service) GetAllByCategory(ctx context.Context, userID int, category string) ([]Transaction, error) {
	return s.repo.GetFiltered(ctx, Filter{
		UserID:   userID,
		Category: &category,
	})
}

func (s *Service) DeleteTransaction(ctx context.Context, userID, id int) error {
	return s.repo.Delete(ctx, userID, id)
}

func (s *Service) UpdateTransaction(ctx context.Context, userID int, t Transaction) (Transaction, error) {
	if err := validate(t); err != nil {
		return Transaction{}, err
	}
	return s.repo.Update(ctx, userID, t)
}

func (s *Service) GetByID(ctx context.Context, userID, id int) (Transaction, error) {
	return s.repo.GetByID(ctx, userID, id)
}
