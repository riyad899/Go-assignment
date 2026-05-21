package user

type service struct {
	repo Repository
}

func NewService(repo Repository) *service {
	return &service{
		repo: repo,
	}
}

func (s *service) CreateUser(user *User) error {
	return s.repo.CreateUser(user)
}
