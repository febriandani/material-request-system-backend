package department

type Service interface {
	GetAll() ([]Department, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetAll() ([]Department, error) {
	return s.repo.FindAll()
}
