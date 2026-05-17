package chat

import (
	"fmt"

	"github.com/DeliriumWaltz/silver-engine/internal/models"
	"github.com/DeliriumWaltz/silver-engine/internal/storage"
)

// Service implements the core chat business logic.
type Service struct {
	store storage.Store
}

// NewService creates a new chat Service.
func NewService(store storage.Store) *Service {
	return &Service{store: store}
}

// CreateRoom creates a new chat room.
func (s *Service) CreateRoom(name string) (*models.Room, error) {
	if name == "" {
		return nil, fmt.Errorf("room name cannot be empty")
	}
	return s.store.CreateRoom(name)
}

// GetRoom retrieves a room by ID.
func (s *Service) GetRoom(id string) (*models.Room, error) {
	return s.store.GetRoom(id)
}

// ListRooms returns all rooms.
func (s *Service) ListRooms() ([]*models.Room, error) {
	return s.store.ListRooms()
}

// SendMessage stores a new message in a room.
func (s *Service) SendMessage(roomID, username, content string) (*models.Message, error) {
	if username == "" {
		return nil, fmt.Errorf("username cannot be empty")
	}
	if content == "" {
		return nil, fmt.Errorf("message content cannot be empty")
	}

	msg := &models.Message{
		RoomID:   roomID,
		Username: username,
		Content:  content,
	}
	if err := s.store.SaveMessage(msg); err != nil {
		return nil, err
	}
	return msg, nil
}

// GetMessages retrieves all messages for a room.
func (s *Service) GetMessages(roomID string) ([]*models.Message, error) {
	return s.store.GetMessages(roomID)
}