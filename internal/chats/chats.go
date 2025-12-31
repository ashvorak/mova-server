package chats

import (
	"slices"

	"github.com/google/uuid"
)

type Chat struct {
	ID      string
	UserIDs []string
}

type Service struct {
	chats map[string]Chat
}

func NewService() *Service {
	return &Service{
		chats: make(map[string]Chat),
	}
}

func (s *Service) Create(userIDs []string) Chat {
	id := uuid.New().String()

	c := Chat{
		ID:      id,
		UserIDs: userIDs,
	}

	s.chats[id] = c
	return c
}

func (s *Service) ListByUser(userID string) []Chat {
	chats := make([]Chat, 0)

	for _, c := range s.chats {
		if slices.Contains(c.UserIDs, userID) {
			chats = append(chats, c)
		}
	}

	return chats
}
