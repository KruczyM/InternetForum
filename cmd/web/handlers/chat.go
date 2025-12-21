package handlers

import (
	"encoding/json"
	"forum/internal/models"
	"net/http"
)

// ChatMessageRequest is the expected payload for posting a chat message
// (can be extended for authentication, etc.)
type ChatMessageRequest struct {
	UserID   string `json:"user_id"`
	Content  string `json:"content"`
	Category string `json:"category"`
}

// ChatMessagesResponse is the response for chat messages by category
// (can be extended for pagination, etc.)
type ChatMessagesResponse struct {
	Messages []models.ChatMessage `json:"messages"`
}

// PostChatMessage handles posting a new chat message
func (h *Handler) PostChatMessage(w http.ResponseWriter, r *http.Request) {
	var req ChatMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.clientError(w, http.StatusBadRequest)
		return
	}
	msg := models.AddChatMessage(req.UserID, req.Content, req.Category)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(msg)
}

// GetChatMessages handles retrieving chat messages for a category
func (h *Handler) GetChatMessages(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	if category == "" {
		h.clientError(w, http.StatusBadRequest)
		return
	}
	msgs := models.GetChatMessages(category)
	resp := ChatMessagesResponse{Messages: msgs}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GetChatCategories returns all available chat categories
func (h *Handler) GetChatCategories(w http.ResponseWriter, r *http.Request) {
	cats := models.GetAllCategories()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cats)
}
