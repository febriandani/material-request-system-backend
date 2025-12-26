package user

type Service interface {
	GetAll() ([]User, error)
	GetApprovers() ([]Approver, error)
}
type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetAll() ([]User, error) {
	return s.repo.FindAll()
}

func (s *service) GetApprovers() ([]Approver, error) {
	return s.repo.FindApprovers()
}
