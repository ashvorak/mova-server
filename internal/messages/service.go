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
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(chatID chats.ID, userID users.ID, text string) (Message, error) {
	m := Message{
		ID:        newID(),
		ChatID:    chatID,
		UserID:    userID,
		Text:      text,
		CreatedAt: time.Now(),
	}

	if err := s.repo.Save(m); err != nil {
		return Message{}, err
	}

	return m, nil
}

func (s *Service) ListByChat(chatID chats.ID) ([]Message, error) {
	return s.repo.ListByChat(chatID)
}

func (s *Service) ListByChatAfter(chatID chats.ID, after ID, limit int) ([]Message, error) {
	messages, err := s.repo.ListByChat(chatID)
	if err != nil {
		return nil, err
	}

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
		return messages, nil
	}

	end := start + limit
	if end > len(messages) {
		end = len(messages)
	}

	result := make([]Message, end-start)
	copy(result, messages[start:end])
	return result, nil
}
