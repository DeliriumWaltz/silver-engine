package models

import "time"

// Message represents a chat message in a room.
type Message struct {
	ID        string    `json:"id"`
	RoomID    string    `json:"room_id"`
	Username  string    `json:"username"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// User represents a connected chat user.
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

// Room represents a chat room.
type Room struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}