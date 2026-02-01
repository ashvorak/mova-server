package messages

import "mova-server/internal/chats"

type Repository interface {
	Save(msg Message) error
	ListByChat(chatID chats.ID) ([]Message, error)
}
