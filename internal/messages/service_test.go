package messages

import (
	"mova-server/internal/chats"
	"mova-server/internal/users"
	"testing"
)

type fakeRepository struct {
	mapByChat map[chats.ID][]Message
}

func (r *fakeRepository) Save(msg Message) error {
	r.mapByChat[msg.ChatID] = append(r.mapByChat[msg.ChatID], msg)
	return nil
}

func (r *fakeRepository) ListByChat(chatID chats.ID) ([]Message, error) {
	return r.mapByChat[chatID], nil
}

func TestService_Create(t *testing.T) {
	repo := &fakeRepository{
		mapByChat: make(map[chats.ID][]Message),
	}
	service := NewService(repo)

	chatsID := chats.NewID()
	userID := users.NewID()

	msg, err := service.Create(chatsID, userID, "Hello, World!")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if msg.ChatID != chatsID {
		t.Errorf("expected chat ID %s, got %s", chatsID, msg.ChatID)
	}

	if msg.UserID != userID {
		t.Errorf("expected user ID %s, got %s", userID, msg.UserID)
	}

	if msg.Text != "Hello, World!" {
		t.Errorf("expected text 'Hello, World!', got '%s'", msg.Text)
	}
}

func TestService_ListByChat(t *testing.T) {
	repo := &fakeRepository{
		mapByChat: make(map[chats.ID][]Message),
	}
	service := NewService(repo)

	chatsID := chats.NewID()
	userID := users.NewID()

	service.Create(chatsID, userID, "First")
	service.Create(chatsID, userID, "Second")
	service.Create(chatsID, userID, "Third")

	messages, err := service.ListByChat(chatsID)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(messages) != 3 {
		t.Errorf("expected 3 messages, got %d", len(messages))
	}

	expectedTexts := []string{"First", "Second", "Third"}
	for i, msg := range messages {
		if msg.Text != expectedTexts[i] {
			t.Errorf("expected %s text, got %s", expectedTexts[i], msg.Text)
		}
	}
}
