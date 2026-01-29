package messages

import "errors"

var (
	ErrMessagesNotFound = errors.New("messages: not found")
	ErrEmptyChatID      = errors.New("messages: empty chat ID")
	ErrEmptyUserID      = errors.New("messages: empty user ID")
	ErrEmptyText        = errors.New("messages: empty text")
)
