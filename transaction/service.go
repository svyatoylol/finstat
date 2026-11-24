package transaction

import (
	"context"
	"errors"
	"fmt"
	"math"
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
	if t.Amount <= 0 || math.IsNaN(t.Amount) || math.IsInf(t.Amount, 0) {
		return ErrInvalidAmount
	}
	if t.Type != Incoming && t.Type != Outgoing {
		return ErrInvalidType
	}
	return nil
}

func (s *Service) GetBalance(ctx context.Context, userID int64) (float64, error) {
	balance, err := s.repo.GetBalance(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("get balance: %w", err)
	}
	return balance, nil
}

func (s *Service) AddTransaction(ctx context.Context, t Transaction) (Transaction, error) {
	if err := validate(t); err != nil {
		return Transaction{}, err
	}
	return s.repo.Add(ctx, t)
}

func (s *Service) GetAllByType(ctx context.Context, userID int64, tType TransactionType) ([]Transaction, error) {
	return s.GetAllByTypePage(ctx, userID, tType, 0, 0)
}

func (s *Service) GetAllPage(ctx context.Context, userID int64, limit int, afterID int64) ([]Transaction, error) {
	return s.repo.GetFiltered(ctx, Filter{
		UserID:  userID,
		Limit:   limit,
		AfterID: afterID,
	})
}

func (s *Service) GetAllByTypePage(ctx context.Context, userID int64, tType TransactionType, limit int, afterID int64) ([]Transaction, error) {
	return s.repo.GetFiltered(ctx, Filter{
		UserID:  userID,
		Type:    &tType,
		Limit:   limit,
		AfterID: afterID,
	})
}

func (s *Service) GetAllByCategory(ctx context.Context, userID int64, category string) ([]Transaction, error) {
	return s.GetAllByCategoryPage(ctx, userID, category, 0, 0)
}

func (s *Service) GetAllByCategoryPage(ctx context.Context, userID int64, category string, limit int, afterID int64) ([]Transaction, error) {
	return s.repo.GetFiltered(ctx, Filter{
		UserID:   userID,
		Category: &category,
		Limit:    limit,
		AfterID:  afterID,
	})
}

func (s *Service) DeleteTransaction(ctx context.Context, userID, id int64) error {
	return s.repo.Delete(ctx, userID, id)
}

func (s *Service) UpdateTransaction(ctx context.Context, userID int64, t Transaction) (Transaction, error) {
	if err := validate(t); err != nil {
		return Transaction{}, err
	}
	return s.repo.Update(ctx, userID, t)
}

func (s *Service) GetByID(ctx context.Context, userID, id int64) (Transaction, error) {
	return s.repo.GetByID(ctx, userID, id)
}
