package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	addr := ":" + port
	srv := &http.Server{Addr: addr, Handler: mux}

	// Graceful shutdown on SIGINT/SIGTERM
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		log.Println("Shutting down...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		srv.Shutdown(ctx)
	}()

	log.Printf("Chat server starting on %s", addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server failed: %v", err)
	}
}