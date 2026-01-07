package chats

import (
	"mova-server/internal/users"
)

type Chat struct {
	ID      ID
	UserIDs []users.ID
}

type Service struct {
	chats       map[ID]Chat
	chatsByUser map[users.ID][]ID
}

func NewService() *Service {
	return &Service{
		chats: make(map[ID]Chat),
	}
}

func (s *Service) Create(userIDs []users.ID) Chat {
	id := newID()

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

	// TODO: Remove
	parsedUserID, err := users.ParseID(userID)
	if err != nil {
		return chats
	}

	for _, chatID := range s.chatsByUser[parsedUserID] {
		chats = append(chats, s.chats[chatID])
	}

	return chats
}
