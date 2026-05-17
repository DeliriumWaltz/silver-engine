package chat

import (
	"testing"

	"github.com/DeliriumWaltz/silver-engine/internal/storage"
)

func setupTestService() *Service {
	return NewService(storage.NewMemoryStore())
}

func TestService_CreateRoom(t *testing.T) {
	s := setupTestService()

	room, err := s.CreateRoom("general")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if room.Name != "general" {
		t.Fatalf("expected name 'general', got %q", room.Name)
	}
}

func TestService_CreateRoom_EmptyName(t *testing.T) {
	s := setupTestService()

	_, err := s.CreateRoom("")
	if err == nil {
		t.Fatal("expected error for empty room name")
	}
}

func TestService_CreateRoom_Duplicate(t *testing.T) {
	s := setupTestService()

	_, err := s.CreateRoom("general")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = s.CreateRoom("general")
	if err == nil {
		t.Fatal("expected error for duplicate room name")
	}
}

func TestService_GetRoom(t *testing.T) {
	s := setupTestService()

	created, err := s.CreateRoom("general")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := s.GetRoom(created.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != created.ID {
		t.Fatalf("expected ID %q, got %q", created.ID, got.ID)
	}
}

func TestService_GetRoom_NotFound(t *testing.T) {
	s := setupTestService()

	_, err := s.GetRoom("nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent room")
	}
}

func TestService_ListRooms(t *testing.T) {
	s := setupTestService()

	rooms, err := s.ListRooms()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rooms) != 0 {
		t.Fatalf("expected 0 rooms, got %d", len(rooms))
	}

	s.CreateRoom("alpha")
	s.CreateRoom("beta")

	rooms, err = s.ListRooms()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rooms) != 2 {
		t.Fatalf("expected 2 rooms, got %d", len(rooms))
	}
}

func TestService_SendMessage(t *testing.T) {
	s := setupTestService()

	room, err := s.CreateRoom("general")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	msg, err := s.SendMessage(room.ID, "alice", "hello")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg.Username != "alice" {
		t.Fatalf("expected username 'alice', got %q", msg.Username)
	}
	if msg.Content != "hello" {
		t.Fatalf("expected content 'hello', got %q", msg.Content)
	}
}

func TestService_SendMessage_EmptyUsername(t *testing.T) {
	s := setupTestService()

	room, err := s.CreateRoom("general")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = s.SendMessage(room.ID, "", "hello")
	if err == nil {
		t.Fatal("expected error for empty username")
	}
}

func TestService_SendMessage_EmptyContent(t *testing.T) {
	s := setupTestService()

	room, err := s.CreateRoom("general")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = s.SendMessage(room.ID, "alice", "")
	if err == nil {
		t.Fatal("expected error for empty content")
	}
}

func TestService_GetMessages(t *testing.T) {
	s := setupTestService()

	room, err := s.CreateRoom("general")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	s.SendMessage(room.ID, "alice", "first")
	s.SendMessage(room.ID, "bob", "second")

	messages, err := s.GetMessages(room.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(messages) != 2 {
		t.Fatalf("expected 2 messages, got %d", len(messages))
	}
}

func TestService_GetMessages_RoomNotFound(t *testing.T) {
	s := setupTestService()

	_, err := s.GetMessages("nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent room")
	}
}