package chats

import (
	"github.com/google/uuid"
)

type Chat struct {
	ID      string
	UserIDs []string
}

type Service struct {
	chats       map[string]Chat     // chatID -> Chat
	chatsByUser map[string][]string // userID -> []chatID
}

func NewService() *Service {
	return &Service{
		chats:       make(map[string]Chat),
		chatsByUser: make(map[string][]string),
	}
}

func (s *Service) Create(userIDs []string) Chat {
	id := uuid.New().String()

	c := Chat{
		ID:      id,
		UserIDs: userIDs,
	}
	s.chats[id] = c

	for _, u := range userIDs {
		s.chatsByUser[u] = append(s.chatsByUser[u], c.ID)
	}

	return c
}

func (s *Service) ListByUser(userID string) []Chat {
	chats := make([]Chat, 0)

	for _, chatID := range s.chatsByUser[userID] {
		chats = append(chats, s.chats[chatID])
	}

	return chats
}
