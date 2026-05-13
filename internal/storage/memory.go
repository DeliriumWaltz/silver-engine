package storage

import (
	"fmt"
	"sync"
	"time"

	"github.com/DeliriumWaltz/silver-engine/internal/models"
	"github.com/google/uuid"
)

// MemoryStore is an in-memory implementation of Store.
// Thread-safe via sync.RWMutex.
type MemoryStore struct {
	mu       sync.RWMutex
	rooms    map[string]*models.Room
	messages map[string][]*models.Message // keyed by roomID
}

// NewMemoryStore creates a new empty MemoryStore.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		rooms:    make(map[string]*models.Room),
		messages: make(map[string][]*models.Message),
	}
}

func (s *MemoryStore) CreateRoom(name string) (*models.Room, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check for duplicate name
	for _, r := range s.rooms {
		if r.Name == name {
			return nil, fmt.Errorf("room with name %q already exists", name)
		}
	}

	room := &models.Room{
		ID:   uuid.NewString(),
		Name: name,
	}
	s.rooms[room.ID] = room
	s.messages[room.ID] = []*models.Message{}
	return room, nil
}

func (s *MemoryStore) GetRoom(id string) (*models.Room, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	room, ok := s.rooms[id]
	if !ok {
		return nil, fmt.Errorf("room %q not found", id)
	}
	return room, nil
}

func (s *MemoryStore) ListRooms() ([]*models.Room, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	rooms := make([]*models.Room, 0, len(s.rooms))
	for _, r := range s.rooms {
		rooms = append(rooms, r)
	}
	return rooms, nil
}

func (s *MemoryStore) SaveMessage(msg *models.Message) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.rooms[msg.RoomID]; !ok {
		return fmt.Errorf("room %q not found", msg.RoomID)
	}

	msg.ID = uuid.NewString()
	msg.CreatedAt = time.Now()
	s.messages[msg.RoomID] = append(s.messages[msg.RoomID], msg)
	return nil
}

func (s *MemoryStore) GetMessages(roomID string) ([]*models.Message, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if _, ok := s.rooms[roomID]; !ok {
		return nil, fmt.Errorf("room %q not found", roomID)
	}

	// Return a copy to avoid external mutation
	msgs := s.messages[roomID]
	result := make([]*models.Message, len(msgs))
	copy(result, msgs)
	return result, nil
}