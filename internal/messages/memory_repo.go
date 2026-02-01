package messages

import "mova-server/internal/chats"

type MemoryRepository struct {
	messages map[chats.ID][]Message
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		messages: make(map[chats.ID][]Message),
	}
}

func (r *MemoryRepository) Save(msg Message) error {
	_, ok := r.messages[msg.ChatID]
	if !ok {
		r.messages[msg.ChatID] = make([]Message, 0)
	}

	if msg.ChatID.IsEmpty() {
		return ErrEmptyChatID
	}

	if msg.UserID.IsEmpty() {
		return ErrEmptyUserID
	}

	if msg.Text == "" {
		return ErrEmptyText
	}

	r.messages[msg.ChatID] = append(r.messages[msg.ChatID], msg)
	return nil
}

func (r *MemoryRepository) ListByChat(chatID chats.ID) ([]Message, error) {
	src, ok := r.messages[chatID]
	if !ok {
		return nil, ErrMessagesNotFound
	}

	dst := make([]Message, len(src))
	copy(dst, src)
	return dst, nil
}
