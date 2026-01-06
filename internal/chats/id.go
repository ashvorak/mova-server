package chats

import "mova-server/internal/shared/id"

type ID id.ID

func (i ID) String() string {
	return string(i)
}

func ParseID(s string) (ID, error) {
	parsed, err := id.Parse(s)
	if err != nil {
		return "", err
	}
	return ID(parsed), nil
}

func newID() ID {
	return ID(id.New())
}
