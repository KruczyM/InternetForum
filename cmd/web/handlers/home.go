package handlers

import (
	"fmt"
	"forum/internal/models"
	"html/template"
	"log"
	"net/http"
)

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	posts, err := models.GetAllPosts(h.DB)
	if err != nil {
		log.Println("DB Error:", err)
		http.Error(w, "Internal Server Error", 500)
		return
	}

	fmt.Fprintf(w, "Welcome to the Lion's Forum!\n\n")
	for _, p := range posts {
		fmt.Fprintf(w, "ID: %d | Title: %s | Content: %s\n", p.ID, p.Title, p.Content)
	}
}

func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		ts, err := template.ParseFiles("ui/html/create.html")
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal Server Error", 500)
			return
		}
		ts.Execute(w, nil)
		return
	}

	if r.Method == http.MethodPost {
		title := r.FormValue("title")
		content := r.FormValue("content")
		userID := r.FormValue("user_id")

		err := models.CreatePost(h.DB, userID, title, content)
		if err != nil {
			log.Println(err)
			http.Error(w, "Error saving post", 500)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
