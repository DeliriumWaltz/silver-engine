package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DeliriumWaltz/silver-engine/internal/chat"
	"github.com/DeliriumWaltz/silver-engine/internal/storage"
)

func setupTestHandler() *Handler {
	store := storage.NewMemoryStore()
	chatService := chat.NewService(store)
	return NewHandler(chatService)
}

func setupTestServer(t *testing.T) (*httptest.Server, *Handler) {
	t.Helper()
	h := setupTestHandler()
	wsHub := NewWSHub(chat.NewService(storage.NewMemoryStore()))
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)
	h.RegisterWSRoute(mux, wsHub)
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)
	return server, h
}

func TestHandler_ListRooms_Empty(t *testing.T) {
	server, _ := setupTestServer(t)

	resp, err := http.Get(server.URL + "/api/rooms")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", resp.StatusCode)
	}

	var rooms []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&rooms); err != nil {
		t.Fatalf("unexpected error decoding response: %v", err)
	}
	if len(rooms) != 0 {
		t.Fatalf("expected 0 rooms, got %d", len(rooms))
	}
}

func TestHandler_CreateRoom(t *testing.T) {
	server, _ := setupTestServer(t)

	body := `{"name":"general"}`
	resp, err := http.Post(server.URL+"/api/rooms", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201 Created, got %d", resp.StatusCode)
	}

	var room map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&room); err != nil {
		t.Fatalf("unexpected error decoding response: %v", err)
	}
	if room["name"] != "general" {
		t.Fatalf("expected name 'general', got %v", room["name"])
	}
	if room["id"] == "" {
		t.Fatal("expected room ID to be non-empty")
	}
}

func TestHandler_CreateRoom_EmptyName(t *testing.T) {
	server, _ := setupTestServer(t)

	body := `{"name":""}`
	resp, err := http.Post(server.URL+"/api/rooms", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request, got %d", resp.StatusCode)
	}
}

func TestHandler_CreateRoom_InvalidJSON(t *testing.T) {
	server, _ := setupTestServer(t)

	resp, err := http.Post(server.URL+"/api/rooms", "application/json", strings.NewReader("not json"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request, got %d", resp.StatusCode)
	}
}

func TestHandler_CreateRoom_Duplicate(t *testing.T) {
	server, _ := setupTestServer(t)

	body := `{"name":"general"}`
	resp, err := http.Post(server.URL+"/api/rooms", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resp.Body.Close()

	// Create duplicate
	resp, err = http.Post(server.URL+"/api/rooms", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request for duplicate, got %d", resp.StatusCode)
	}
}

func TestHandler_GetRoom(t *testing.T) {
	server, _ := setupTestServer(t)

	// Create a room first
	body := `{"name":"test-room"}`
	resp, err := http.Post(server.URL+"/api/rooms", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var created map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&created)
	resp.Body.Close()

	// Get the room
	resp, err = http.Get(server.URL + "/api/rooms/" + created["id"].(string))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", resp.StatusCode)
	}

	var got map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&got)
	if got["name"] != "test-room" {
		t.Fatalf("expected name 'test-room', got %v", got["name"])
	}
}

func TestHandler_GetRoom_NotFound(t *testing.T) {
	server, _ := setupTestServer(t)

	resp, err := http.Get(server.URL + "/api/rooms/nonexistent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404 Not Found, got %d", resp.StatusCode)
	}
}

func TestHandler_SendAndGetMessages(t *testing.T) {
	server, _ := setupTestServer(t)

	// Create a room
	body := `{"name":"chat-room"}`
	resp, err := http.Post(server.URL+"/api/rooms", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var room map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&room)
	resp.Body.Close()

	roomID := room["id"].(string)

	// Send a message
	msgBody := `{"username":"alice","content":"hello world"}`
	resp, err = http.Post(server.URL+"/api/rooms/"+roomID+"/messages", "application/json", strings.NewReader(msgBody))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201 Created, got %d", resp.StatusCode)
	}

	var sentMsg map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&sentMsg)
	if sentMsg["username"] != "alice" {
		t.Fatalf("expected username 'alice', got %v", sentMsg["username"])
	}
	if sentMsg["content"] != "hello world" {
		t.Fatalf("expected content 'hello world', got %v", sentMsg["content"])
	}

	// Get messages
	resp, err = http.Get(server.URL + "/api/rooms/" + roomID + "/messages")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", resp.StatusCode)
	}

	var messages []map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&messages)
	if len(messages) != 1 {
		t.Fatalf("expected 1 message, got %d", len(messages))
	}
	if messages[0]["content"] != "hello world" {
		t.Fatalf("expected content 'hello world', got %v", messages[0]["content"])
	}
}

func TestHandler_SendMessage_EmptyUsername(t *testing.T) {
	server, _ := setupTestServer(t)

	// Create a room
	resp, _ := http.Post(server.URL+"/api/rooms", "application/json", strings.NewReader(`{"name":"r"}`))
	var room map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&room)
	resp.Body.Close()

	msgBody := `{"username":"","content":"hello"}`
	resp, err := http.Post(server.URL+"/api/rooms/"+room["id"].(string)+"/messages", "application/json", strings.NewReader(msgBody))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request for empty username, got %d", resp.StatusCode)
	}
}

func TestHandler_SendMessage_EmptyContent(t *testing.T) {
	server, _ := setupTestServer(t)

	// Create a room
	resp, _ := http.Post(server.URL+"/api/rooms", "application/json", strings.NewReader(`{"name":"r"}`))
	var room map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&room)
	resp.Body.Close()

	msgBody := `{"username":"alice","content":""}`
	resp, err := http.Post(server.URL+"/api/rooms/"+room["id"].(string)+"/messages", "application/json", strings.NewReader(msgBody))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request for empty content, got %d", resp.StatusCode)
	}
}

func TestHandler_SendMessage_RoomNotFound(t *testing.T) {
	server, _ := setupTestServer(t)

	msgBody := `{"username":"alice","content":"hello"}`
	resp, err := http.Post(server.URL+"/api/rooms/nonexistent/messages", "application/json", strings.NewReader(msgBody))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request for nonexistent room, got %d", resp.StatusCode)
	}
}

func TestHandler_GetMessages_RoomNotFound(t *testing.T) {
	server, _ := setupTestServer(t)

	resp, err := http.Get(server.URL + "/api/rooms/nonexistent/messages")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404 Not Found, got %d", resp.StatusCode)
	}
}

func TestHandler_ListRooms_AfterCreate(t *testing.T) {
	server, _ := setupTestServer(t)

	// Create two rooms
	http.Post(server.URL+"/api/rooms", "application/json", strings.NewReader(`{"name":"alpha"}`))
	http.Post(server.URL+"/api/rooms", "application/json", strings.NewReader(`{"name":"beta"}`))

	resp, err := http.Get(server.URL + "/api/rooms")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	var rooms []map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&rooms)
	if len(rooms) != 2 {
		t.Fatalf("expected 2 rooms, got %d", len(rooms))
	}
}

func TestHandler_Integration_FullFlow(t *testing.T) {
	// Full end-to-end integration test
	server, _ := setupTestServer(t)

	// 1. Create room
	resp, _ := http.Post(server.URL+"/api/rooms", "application/json", bytes.NewReader([]byte(`{"name":"integration-test"}`)))
	var room map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&room)
	resp.Body.Close()
	roomID := room["id"].(string)

	// 2. Send messages
	for i, user := range []string{"alice", "bob", "charlie"} {
		body, _ := json.Marshal(map[string]string{"username": user, "content": "message " + string(rune('0'+i))})
		http.Post(server.URL+"/api/rooms/"+roomID+"/messages", "application/json", bytes.NewReader(body))
	}

	// 3. Get messages
	resp, _ = http.Get(server.URL + "/api/rooms/" + roomID + "/messages")
	var messages []map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&messages)
	resp.Body.Close()

	if len(messages) != 3 {
		t.Fatalf("expected 3 messages, got %d", len(messages))
	}

	// 4. List rooms
	resp, _ = http.Get(server.URL + "/api/rooms")
	var rooms []map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&rooms)
	resp.Body.Close()

	if len(rooms) < 1 {
		t.Fatal("expected at least 1 room")
	}
}

func TestWSHub_ServeWS_NoRoomID(t *testing.T) {
	handler := setupTestHandler()
	wsHub := NewWSHub(chat.NewService(storage.NewMemoryStore()))
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)
	handler.RegisterWSRoute(mux, wsHub)
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	resp, err := http.Get(server.URL + "/api/rooms//ws")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	// With our route pattern, empty roomID won't match, so we'll get 404
	// which means the route didn't match — acceptable for now
	if resp.StatusCode != http.StatusNotFound && resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 404 or 400, got %d", resp.StatusCode)
	}
}

func TestHandler_POST_Room_WithSpaces(t *testing.T) {
	server, _ := setupTestServer(t)

	body := `{"name":"room with spaces"}`
	resp, err := http.Post(server.URL+"/api/rooms", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201 Created, got %d", resp.StatusCode)
	}
}
