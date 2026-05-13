package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/DeliriumWaltz/silver-engine/internal/api"
	"github.com/DeliriumWaltz/silver-engine/internal/chat"
	"github.com/DeliriumWaltz/silver-engine/internal/storage"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Use in-memory storage (swap with DB-backed store later)
	store := storage.NewMemoryStore()
	chatService := chat.NewService(store)
	handler := api.NewHandler(chatService)
	wsHub := api.NewWSHub(chatService)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)
	handler.RegisterWSRoute(mux, wsHub)

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Chat server starting on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}