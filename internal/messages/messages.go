package messages

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID        string
	ChatID    string
	UserID    string
	Text      string
	CreatedAt time.Time
}

type Service struct {
	messages map[string][]Message // chatID -> messages
}

func NewService() *Service {
	return &Service{
		messages: make(map[string][]Message),
	}
}

func (s *Service) Create(chatID string, userID string, text string) Message {
	m := Message{
		ID:        uuid.New().String(),
		ChatID:    chatID,
		UserID:    userID,
		Text:      text,
		CreatedAt: time.Now(),
	}

	s.messages[chatID] = append(s.messages[chatID], m)
	return m
}

func (s *Service) ListByChat(chatID string) []Message {
	src := s.messages[chatID]

	dst := make([]Message, len(src))
	copy(dst, src)

	return dst
}
