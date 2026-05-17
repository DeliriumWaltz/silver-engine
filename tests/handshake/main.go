//go:build ignore

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func main() {
	baseURL := "http://localhost:8080"

	// 1. List rooms (empty)
	fmt.Println("=== List rooms (expect empty) ===")
	resp, err := http.Get(baseURL + "/api/rooms")
	if err != nil {
		panic(err)
	}
	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("Status: %d, Body: %s\n\n", resp.StatusCode, body)
	resp.Body.Close()

	// 2. Create room
	fmt.Println("=== Create room ===")
	resp, err = http.Post(baseURL+"/api/rooms", "application/json", strings.NewReader(`{"name":"general"}`))
	if err != nil {
		panic(err)
	}
	body, _ = io.ReadAll(resp.Body)
	fmt.Printf("Status: %d, Body: %s\n\n", resp.StatusCode, body)
	resp.Body.Close()

	var room map[string]interface{}
	json.Unmarshal(body, &room)
	roomID := room["id"].(string)

	// 3. Send message
	fmt.Println("=== Send message ===")
	msgBody := `{"username":"alice","content":"hello world"}`
	resp, err = http.Post(baseURL+"/api/rooms/"+roomID+"/messages", "application/json", strings.NewReader(msgBody))
	if err != nil {
		panic(err)
	}
	body, _ = io.ReadAll(resp.Body)
	fmt.Printf("Status: %d, Body: %s\n\n", resp.StatusCode, body)
	resp.Body.Close()

	// 4. Get messages
	fmt.Println("=== Get messages ===")
	resp, err = http.Get(baseURL + "/api/rooms/" + roomID + "/messages")
	if err != nil {
		panic(err)
	}
	body, _ = io.ReadAll(resp.Body)
	fmt.Printf("Status: %d, Body: %s\n\n", resp.StatusCode, body)
	resp.Body.Close()

	// 5. List rooms (should have 1)
	fmt.Println("=== List rooms (expect 1) ===")
	resp, err = http.Get(baseURL + "/api/rooms")
	if err != nil {
		panic(err)
	}
	body, _ = io.ReadAll(resp.Body)
	fmt.Printf("Status: %d, Body: %s\n", resp.StatusCode, body)
	resp.Body.Close()

	fmt.Println("All smoke tests passed!")
}