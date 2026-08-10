package transaction

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
		if t.Type == "incoming" {
			balance += t.Amount
		}
		if t.Type == "outgoing" {
			balance -= t.Amount
		}
	}
	return balance
}
