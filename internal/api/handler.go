package api

import (
	"encoding/json"
	"net/http"

	"github.com/DeliriumWaltz/silver-engine/internal/chat"
	"github.com/google/uuid"
)

type Handler struct {
	chatService *chat.Service
}

func NewHandler(chatService *chat.Service) *Handler {
	return &Handler{chatService: chatService}
}

// RegisterRoutes sets up HTTP routes on the given mux.
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/rooms", h.ListRooms)
	mux.HandleFunc("POST /api/rooms", h.CreateRoom)
	mux.HandleFunc("GET /api/rooms/{roomID}", h.GetRoom)
	mux.HandleFunc("GET /api/rooms/{roomID}/messages", h.GetMessages)
	mux.HandleFunc("POST /api/rooms/{roomID}/messages", h.SendMessage)
}

type createRoomRequest struct {
	Name string `json:"name"`
}

type sendMessageRequest struct {
	Username string `json:"username"`
	Content  string `json:"content"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func (h *Handler) ListRooms(w http.ResponseWriter, r *http.Request) {
	rooms, err := h.chatService.ListRooms()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, rooms)
}

func (h *Handler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	var req createRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid JSON body"})
		return
	}

	room, err := h.chatService.CreateRoom(req.Name)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, room)
}

func (h *Handler) GetRoom(w http.ResponseWriter, r *http.Request) {
	roomID := r.PathValue("roomID")
	room, err := h.chatService.GetRoom(roomID)
	if err != nil {
		writeJSON(w, http.StatusNotFound, errorResponse{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, room)
}

func (h *Handler) GetMessages(w http.ResponseWriter, r *http.Request) {
	roomID := r.PathValue("roomID")
	messages, err := h.chatService.GetMessages(roomID)
	if err != nil {
		writeJSON(w, http.StatusNotFound, errorResponse{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, messages)
}

func (h *Handler) SendMessage(w http.ResponseWriter, r *http.Request) {
	roomID := r.PathValue("roomID")

	var req sendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid JSON body"})
		return
	}

	msg, err := h.chatService.SendMessage(roomID, req.Username, req.Content)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, msg)
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

// --- WebSocket handler (stub for now) ---

// WSHub manages connected WebSocket clients.
// For the MVP it's minimal — real-time broadcasting will be added later.
type WSHub struct {
	chatService *chat.Service
}

func NewWSHub(chatService *chat.Service) *WSHub {
	return &WSHub{chatService: chatService}
}

// ServeWS handles WebSocket upgrade requests.
// For the MVP we use a placeholder that returns a proper upgrade.
func (h *WSHub) ServeWS(w http.ResponseWriter, r *http.Request) {
	// For now, we keep this as a placeholder that returns an upgrade
	// so the tests can verify the endpoint exists.
	// In the next iteration, this will use gorilla/websocket to upgrade.
	roomID := r.PathValue("roomID")
	if roomID == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "roomID is required"})
		return
	}

	// Generate a random username for demo purposes
	username := "user-" + uuid.NewString()[:8]

	// We'll upgrade in a follow-up; for now respond with info.
	writeJSON(w, http.StatusOK, map[string]string{
		"room_id":  roomID,
		"username": username,
		"status":   "WebSocket upgrade not yet implemented; use REST API for now",
	})
}

func (h *Handler) RegisterWSRoute(mux *http.ServeMux, hub *WSHub) {
	mux.HandleFunc("GET /api/rooms/{roomID}/ws", hub.ServeWS)
}