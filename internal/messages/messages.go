package messages

import (
	"mova-server/internal/chats"
	"mova-server/internal/users"
	"time"
)

type Message struct {
	ID        ID
	ChatID    chats.ID
	UserID    users.ID
	Text      string
	CreatedAt time.Time
}

type Service struct {
	messages map[chats.ID][]Message
}

func NewService() *Service {
	return &Service{
		messages: make(map[chats.ID][]Message),
	}
}

func (s *Service) Create(chatID chats.ID, userID users.ID, text string) Message {
	m := Message{
		ID:        newID(),
		ChatID:    chatID,
		UserID:    userID,
		Text:      text,
		CreatedAt: time.Now(),
	}

	s.messages[chatID] = append(s.messages[chatID], m)
	return m
}

func (s *Service) ListByChat(chatID chats.ID) []Message {
	src := s.messages[chatID]

	dst := make([]Message, len(src))
	copy(dst, src)

	return dst
}
