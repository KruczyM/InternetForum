package models

import (
	"sync"
	"time"
)

type ChatMessage struct {
	ID        int
	UserID    string
	Content   string
	Category  string // Added for filtering by category
	CreatedAt time.Time
}

// In-memory chat storage: map[category][]ChatMessage
var (
	chatStorage = make(map[string][]ChatMessage)
	chatMutex   sync.RWMutex
)

// adds a new chat message to the in-memory storage
func AddChatMessage(userID, content, category string) ChatMessage {
	msg := ChatMessage{
		ID:        int(time.Now().UnixNano()), // Simple unique ID
		UserID:    userID,
		Content:   content,
		Category:  category,
		CreatedAt: time.Now(),
	}
	chatMutex.Lock()
	chatStorage[category] = append(chatStorage[category], msg)
	chatMutex.Unlock()
	return msg
}

// returns all chat messages for a given category
func GetChatMessages(category string) []ChatMessage {
	chatMutex.RLock()
	defer chatMutex.RUnlock()
	return append([]ChatMessage(nil), chatStorage[category]...)
}

// returns a list of all chat categories
func GetAllCategories() []string {
	chatMutex.RLock()
	defer chatMutex.RUnlock()
	cats := make([]string, 0, len(chatStorage))
	for cat := range chatStorage {
		cats = append(cats, cat)
	}
	return cats
}
