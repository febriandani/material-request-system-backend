package user

type Service interface {
	GetAll() ([]User, error)
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
