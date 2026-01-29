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

func (s *Service) Create(userIDs []users.ID) (Chat, error) {
	if len(userIDs) == 0 {
		return Chat{}, ErrEmptyUserIDs
	}

	id := newID()
	c := Chat{
		ID:      id,
		UserIDs: userIDs,
	}
	s.chats[id] = c

	for _, u := range userIDs {
		s.chatsByUser[u] = append(s.chatsByUser[u], c.ID)
	}

	return c, nil
}

func (s *Service) ListByUser(userID users.ID) ([]Chat, error) {
	if _, ok := s.chatsByUser[userID]; !ok {
		return []Chat{}, ErrChatNotFound
	}

	chats := make([]Chat, 0)
	for _, chatID := range s.chatsByUser[userID] {
		chats = append(chats, s.chats[chatID])
	}

	return chats, nil
}
