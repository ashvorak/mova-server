package chats

import "errors"

var (
	ErrChatNotFound = errors.New("chat not found")
	ErrEmptyText    = errors.New("empty message text")
	ErrEmptyUserIDs = errors.New("empty user IDs")
)
