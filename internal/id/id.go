package id

import (
	"errors"

	"github.com/google/uuid"
)

type ID string

func New() ID {
	return ID(uuid.New().String())
}

func Parse(s string) (ID, error) {
	if s == "" {
		return "", errors.New("empty id")
	}

	_, err := uuid.Parse(s)
	if err != nil {
		return "", err
	}

	return ID(s), nil
}

func (id ID) String() string {
	return string(id)
}
