package users

type User struct {
	ID   ID
	Name string
}

type Service struct {
	users map[ID]User
}

func NewService() *Service {
	return &Service{
		users: make(map[ID]User),
	}
}

func (s *Service) Create(name string) (User, error) {
	if name == "" {
		return User{}, ErrUserInvalidName
	}

	id := NewID()

	user := User{
		ID:   id,
		Name: name,
	}

	s.users[id] = user
	return user, nil
}

func (s *Service) List() []User {
	l := make([]User, 0, len(s.users))
	for _, user := range s.users {
		l = append(l, user)
	}

	return l
}
