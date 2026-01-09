package messages

import (
	"mova-server/internal/chats"
	"mova-server/internal/users"
	"slices"
	"time"
)

const (
	defaultListLimit = 50
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

func (s *Service) ListByChatAfter(chatID chats.ID, after ID, limit int) []Message {
	messages := s.messages[chatID]

	start := 0
	if after != "" {
		if idx := slices.IndexFunc(messages, func(m Message) bool {
			return m.ID == after
		}); idx != -1 {
			start = idx + 1
		}
	}

	if limit <= 0 {
		limit = defaultListLimit
	}

	if start >= len(messages) {
		return nil
	}

	end := start + limit
	if end > len(messages) {
		end = len(messages)
	}

	result := make([]Message, end-start)
	copy(result, messages[start:end])
	return result
}
