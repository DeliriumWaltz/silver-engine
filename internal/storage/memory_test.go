package storage

import (
	"testing"

	"github.com/DeliriumWaltz/silver-engine/internal/models"
)

func TestMemoryStore_CreateRoom(t *testing.T) {
	s := NewMemoryStore()

	room, err := s.CreateRoom("general")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if room.ID == "" {
		t.Fatal("expected room ID to be non-empty")
	}
	if room.Name != "general" {
		t.Fatalf("expected room name 'general', got %q", room.Name)
	}
}

func TestMemoryStore_CreateRoom_Duplicate(t *testing.T) {
	s := NewMemoryStore()

	_, err := s.CreateRoom("general")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = s.CreateRoom("general")
	if err == nil {
		t.Fatal("expected error for duplicate room name")
	}
}

func TestMemoryStore_GetRoom_NotFound(t *testing.T) {
	s := NewMemoryStore()

	_, err := s.GetRoom("nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent room")
	}
}

func TestMemoryStore_GetRoom(t *testing.T) {
	s := NewMemoryStore()

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

func TestMemoryStore_ListRooms(t *testing.T) {
	s := NewMemoryStore()

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

func TestMemoryStore_SaveAndGetMessages(t *testing.T) {
	s := NewMemoryStore()

	room, err := s.CreateRoom("general")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	msg := &models.Message{
		RoomID:   room.ID,
		Username: "alice",
		Content:  "hello",
	}
	if err := s.SaveMessage(msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg.ID == "" {
		t.Fatal("expected message ID to be set after save")
	}

	msg2 := &models.Message{
		RoomID:   room.ID,
		Username: "bob",
		Content:  "hi there",
	}
	if err := s.SaveMessage(msg2); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	messages, err := s.GetMessages(room.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(messages) != 2 {
		t.Fatalf("expected 2 messages, got %d", len(messages))
	}
	if messages[0].Username != "alice" {
		t.Fatalf("expected first message from 'alice', got %q", messages[0].Username)
	}
	if messages[1].Username != "bob" {
		t.Fatalf("expected second message from 'bob', got %q", messages[1].Username)
	}
}

func TestMemoryStore_SaveMessage_RoomNotFound(t *testing.T) {
	s := NewMemoryStore()

	msg := &models.Message{
		RoomID:   "nonexistent",
		Username: "alice",
		Content:  "hello",
	}
	if err := s.SaveMessage(msg); err == nil {
		t.Fatal("expected error for nonexistent room")
	}
}

func TestMemoryStore_GetMessages_RoomNotFound(t *testing.T) {
	s := NewMemoryStore()

	_, err := s.GetMessages("nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent room")
	}
}