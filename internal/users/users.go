package users

import (
	"github.com/google/uuid"
)

type User struct {
	ID   string
	Name string
}

type Service struct {
	users map[string]User
}

func NewService() *Service {
	return &Service{
		users: make(map[string]User),
	}
}

func (s *Service) Create(name string) User {
	id := uuid.New().String()

	user := User{
		ID:   id,
		Name: name,
	}

	s.users[id] = user
	return user
}

func (s *Service) List() []User {
	l := make([]User, 0, len(s.users))
	for _, user := range s.users {
		l = append(l, user)
	}

	return l
}
