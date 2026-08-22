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

func (s *Service) GetBalance(userID int) float64 {
	var balance float64

	for _, t := range s.repo.GetFiltered(Filter{UserID: userID}) {
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

func (s *Service) GetAllByType(tType TransactionType, userID int) []Transaction {
	filter := Filter{
		UserID: userID,
		Type:   &tType,
	}
	return s.repo.GetFiltered(filter)
}

func (s *Service) GetAllByCategory(category string, userID int) []Transaction {
	filter := Filter{
		UserID:   userID,
		Category: &category,
	}
	return s.repo.GetFiltered(filter)
}

func (s *Service) DeleteTransaction(userID int, id int) error {
	return s.repo.Delete(userID, id)
}

func (s *Service) UpdateTransaction(userID int, t Transaction) (Transaction, error) {
	if t.Amount <= 0 {
		return Transaction{}, errors.New("amount must be positive")
	}
	if t.Type != Incoming && t.Type != Outgoing {
		return Transaction{}, errors.New("invalid transaction type")
	}
	return s.repo.Update(userID, t)
}

func (s *Service) GetByID(userID int, id int) (Transaction, bool) {
	return s.repo.GetByID(userID, id)
}
