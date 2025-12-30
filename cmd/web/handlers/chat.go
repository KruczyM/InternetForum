package handlers

import (
	"forum/internal/models"
	"net/http"
	"strings"
)


type ChatData struct {
	Categories      []string
	CurrentCategory string
	Messages        []models.ChatMessage
}


// ChatHandler serves the chat page and handles form submissions
func (h *Handler) ChatHandler(w http.ResponseWriter, r *http.Request) {
       // Get all categories
       categories := models.GetAllCategories()
       var currentCategory string
       if len(categories) > 0 {
	       currentCategory = categories[0]
       }
       // Allow category selection via GET or POST
       if cat := r.FormValue("category"); cat != "" {
	       currentCategory = cat
       }

       // Handle POST (send message)
       if r.Method == http.MethodPost {
	       userID := strings.TrimSpace(r.FormValue("user_id"))
	       content := strings.TrimSpace(r.FormValue("content"))
	       if userID != "" && content != "" && currentCategory != "" {
		       models.AddChatMessage(userID, content, currentCategory)
	       }
	       // Redirect to GET to avoid resubmission
	       http.Redirect(w, r, "/chat/?category="+currentCategory, http.StatusSeeOther)
	       return
       }

       // Load messages for the selected category
       messages := models.GetChatMessages(currentCategory)

       data := &ChatData{
	       Categories:      categories,
	       CurrentCategory: currentCategory,
	       Messages:        messages,
       }
       // Render the chat template
       h.render(w, http.StatusOK, "chat.html", &templateData{AnyData: map[string]any{"Chat": data}})
}
