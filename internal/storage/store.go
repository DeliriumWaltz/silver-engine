package storage

import "github.com/DeliriumWaltz/silver-engine/internal/models"

// Store defines the interface for data storage.
// Implementations can be in-memory, PostgreSQL, etc.
type Store interface {
	// CreateRoom creates a new chat room.
	CreateRoom(name string) (*models.Room, error)
	// GetRoom retrieves a room by ID.
	GetRoom(id string) (*models.Room, error)
	// ListRooms returns all rooms.
	ListRooms() ([]*models.Room, error)

	// SaveMessage stores a message in a room.
	SaveMessage(msg *models.Message) error
	// GetMessages retrieves all messages for a room in chronological order.
	GetMessages(roomID string) ([]*models.Message, error)
}