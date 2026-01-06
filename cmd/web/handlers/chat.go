package handlers

import (
	"forum/internal/models"
	"net/http"
	"strings"
)


type ChatData struct {
	Messages []models.ChatMessage
}


// ChatHandler serves the chat page and handles form submissions
func (h *Handler) ChatHandler(w http.ResponseWriter, r *http.Request) {
       data := h.newTemplateData(w,r)
       // Handle POST (send message)
       if r.Method == http.MethodPost {
	       userID := data.AuthenticatedUser
	       content := strings.TrimSpace(r.FormValue("content"))
		   if userID == ""{
			h.setFlash(w, "error", "Please log in to perform this action. Not registered yet? <a href='/auth/register'  class='back-link' style='display: inline-block;'>Register</a>")
 			http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		   }
	       if userID != "" && content != "" {
		       // Fetch first name for chat display
		       user, err := models.GetUserByID(h.DB, userID)
		       name := userID
		       if err == nil && user != nil && user.FirstName != "" {
			       name = user.FirstName
		       }
		       models.AddChatMessage(name, content, "")
	       }
	       // Redirect to GET to avoid resubmission
	       http.Redirect(w, r, "/chat/", http.StatusSeeOther)
	       return
       }

       // Load all messages (no categories)
       messages := models.GetChatMessages("")

       chatData := &ChatData{
	       Messages: messages,
       }
       data.AnyData["Chat"] = chatData
       // Render the chat template
       h.render(w, http.StatusOK, "chat.html", data)
}
