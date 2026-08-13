package transaction

import (
	"errors"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetBalance() float64 {
	var balance float64

	for _, t := range s.repo.GetAll() {
		if t.Type == Incoming {
			balance += t.Amount
		}
		if t.Type == Outgoing {
			balance -= t.Amount
		}
	}
	return balance
}

func (s *Service) AddTransaction(t Transaction) (Transaction, error) {
	if t.Amount <= 0 {
		return Transaction{}, errors.New("amount must be positive")
	}
	if t.Type != Incoming && t.Type != Outgoing {
		return Transaction{}, errors.New("invalid transaction type")
	}

	saved := s.repo.Add(t)
	return saved, nil
}

func (s *Service) GetAllByType(tType TransactionType) []Transaction {
	filter := Filter{
		Type: &tType,
	}
	return s.repo.GetFiltered(filter)
}

func (s *Service) GetAllByCategory(category string) []Transaction {
	filter := Filter{
		Category: &category,
	}
	return s.repo.GetFiltered(filter)
}
