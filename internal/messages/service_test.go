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

	chatID := chats.NewID()
	userID := users.NewID()

	msg, err := service.Create(chatID, userID, "Hello, World!")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if msg.ID == "" {
		t.Errorf("expected message ID to be set")
	}

	if msg.ChatID != chatID {
		t.Errorf("expected chat ID %s, got %s", chatID, msg.ChatID)
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

	if _, err := service.Create(chatsID, userID, "First"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := service.Create(chatsID, userID, "Second"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := service.Create(chatsID, userID, "Third"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

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

func TestService_ListByChatAfter(t *testing.T) {
	repo := &fakeRepository{
		mapByChat: make(map[chats.ID][]Message),
	}
	service := NewService(repo)

	chatID := chats.NewID()
	userID := users.NewID()

	var msg1, msg2, msg3, msg4 Message
	var err error
	if msg1, err = service.Create(chatID, userID, "First"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg2, err = service.Create(chatID, userID, "Second"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg3, err = service.Create(chatID, userID, "Third"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg4, err = service.Create(chatID, userID, "Fourth"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tests := []struct {
		name     string
		after    ID
		limit    int
		expected []Message
	}{
		{
			name:     "no after, default limit",
			after:    "",
			limit:    0,
			expected: []Message{msg1, msg2, msg3, msg4},
		},
		{
			name:     "after second message",
			after:    msg2.ID,
			limit:    0,
			expected: []Message{msg3, msg4},
		},
		{
			name:     "after second message with limit",
			after:    msg2.ID,
			limit:    1,
			expected: []Message{msg3},
		},
		{
			name:     "after last message",
			after:    msg4.ID,
			limit:    10,
			expected: nil,
		},
		{
			name:     "after not existing message",
			after:    ID("non-existing"),
			limit:    0,
			expected: []Message{msg1, msg2, msg3, msg4},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.ListByChatAfter(chatID, tt.after, tt.limit)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(result) != len(tt.expected) {
				t.Fatalf("expected %d messages, got %d", len(tt.expected), len(result))
			}

			for i := range result {
				if result[i].ID != tt.expected[i].ID {
					t.Errorf(
						"at index %d: expected message ID %s, got %s",
						i,
						tt.expected[i].ID,
						result[i].ID,
					)
				}
			}
		})
	}
}
