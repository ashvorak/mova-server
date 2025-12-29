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
