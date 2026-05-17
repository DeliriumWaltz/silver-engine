// Package integration contains end-to-end tests that compile and run the
// chat-server binary as a subprocess and send real HTTP requests to it.
package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

const (
	binaryPath = "../../bin/chat-server"
	testPort   = "19999"
)

// startServer builds and starts the chat-server binary, returning a function to stop it.
func startServer(t *testing.T) string {
	t.Helper()

	// Build the binary if it doesn't exist
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		build := exec.Command("go", "build", "-o", binaryPath, "../../cmd/chat-server")
		build.Stderr = os.Stderr
		if err := build.Run(); err != nil {
			t.Fatalf("failed to build server binary: %v", err)
		}
	}

	addr := fmt.Sprintf("localhost:%s", testPort)
	cmd := exec.Command(binaryPath)
	cmd.Env = append(os.Environ(), "PORT="+testPort)
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		t.Fatalf("failed to start server: %v", err)
	}

	// Cleanup: kill the server when the test finishes
	t.Cleanup(func() {
		cmd.Process.Kill()
		cmd.Wait()
	})

	// Wait for server to be ready
	baseURL := fmt.Sprintf("http://%s", addr)
	if err := waitForServer(baseURL, 5*time.Second); err != nil {
		t.Fatalf("server didn't start in time: %v", err)
	}

	return baseURL
}

// waitForServer polls the server until it responds or times out.
func waitForServer(baseURL string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		resp, err := http.Get(baseURL + "/api/rooms")
		if err == nil {
			resp.Body.Close()
			return nil
		}
		time.Sleep(200 * time.Millisecond)
	}
	return fmt.Errorf("server not reachable at %s", baseURL)
}

func TestBinary_CreateAndListRooms(t *testing.T) {
	baseURL := startServer(t)

	// Create a room
	body := `{"name":"integration-room"}`
	resp, err := http.Post(baseURL+"/api/rooms", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}

	var room map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&room); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if room["name"] != "integration-room" {
		t.Fatalf("expected name 'integration-room', got %v", room["name"])
	}

	roomID := room["id"].(string)

	// List rooms
	resp, err = http.Get(baseURL + "/api/rooms")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	var rooms []map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&rooms)

	found := false
	for _, r := range rooms {
		if r["id"] == roomID {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("created room not found in list")
	}
}

func TestBinary_SendAndGetMessages(t *testing.T) {
	baseURL := startServer(t)

	// Create room
	resp, _ := http.Post(baseURL+"/api/rooms", "application/json", strings.NewReader(`{"name":"msg-test"}`))
	var room map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&room)
	resp.Body.Close()
	roomID := room["id"].(string)

	// Send a message
	msgBody := `{"username":"tester","content":"binary integration test"}`
	resp, err := http.Post(baseURL+"/api/rooms/"+roomID+"/messages", "application/json", strings.NewReader(msgBody))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}

	// Get messages
	resp, err = http.Get(baseURL + "/api/rooms/" + roomID + "/messages")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	var messages []map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&messages)

	if len(messages) != 1 {
		t.Fatalf("expected 1 message, got %d", len(messages))
	}
	if messages[0]["content"] != "binary integration test" {
		t.Fatalf("expected content 'binary integration test', got %v", messages[0]["content"])
	}
}

func TestBinary_MultipleUsers(t *testing.T) {
	baseURL := startServer(t)

	// Create room
	resp, _ := http.Post(baseURL+"/api/rooms", "application/json", strings.NewReader(`{"name":"multi-user"}`))
	var room map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&room)
	resp.Body.Close()
	roomID := room["id"].(string)

	users := []string{"alice", "bob", "charlie"}
	for _, user := range users {
		body := fmt.Sprintf(`{"username":"%s","content":"hello from %s"}`, user, user)
		resp, err := http.Post(baseURL+"/api/rooms/"+roomID+"/messages", "application/json", strings.NewReader(body))
		if err != nil {
			t.Fatalf("unexpected error for %s: %v", user, err)
		}
		resp.Body.Close()
	}

	// Verify all messages are stored
	resp, _ = http.Get(baseURL + "/api/rooms/" + roomID + "/messages")
	var messages []map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&messages)
	resp.Body.Close()

	if len(messages) != 3 {
		t.Fatalf("expected 3 messages, got %d", len(messages))
	}

	for i, user := range users {
		if messages[i]["username"] != user {
			t.Fatalf("expected message %d from %s, got %v", i, user, messages[i]["username"])
		}
	}
}

func TestBinary_RoomNotFound(t *testing.T) {
	baseURL := startServer(t)

	resp, err := http.Get(baseURL + "/api/rooms/nonexistent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404 for nonexistent room, got %d", resp.StatusCode)
	}
}
