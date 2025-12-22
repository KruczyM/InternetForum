package handlers

import (
	"fmt"
	"forum/internal/models"
	"html/template"
	"net/http"
	"strconv"
)

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {

	postsModel := &models.PostModel{DB: h.DB}

	posts, err := postsModel.GetAllPosts()
	if err != nil {
		h.ErrorLog.Println("DB Error:", err)
		h.serverError(w, err)
		return
	}

	data := struct {
		Posts []models.PostView
	}{
		Posts: posts,
	}

	ts, err := template.ParseFiles("ui/html/home.html")
	if err != nil {
		h.ErrorLog.Println("Template Error:", err)
		h.serverError(w, err)
		return
	}
	ts.Execute(w, data)
}

func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		ts, err := template.ParseFiles("ui/html/create.html")
		if err != nil {
			h.ErrorLog.Println("Template Error:", err)
			h.serverError(w, err)
			return
		}
		ts.Execute(w, nil)
		return
	}

	if r.Method == http.MethodPost {

		if r.Method == http.MethodPost {
			err := r.ParseForm()
			if err != nil {
				h.serverError(w, err)
				return
			}

			title := r.FormValue("title")
			content := r.FormValue("content")
			category := r.FormValue("category")
			userID := "1"

			var bookID *int

			rawBookID := r.FormValue("book_id")

			if category == "book" && rawBookID != "" {
				id, err := strconv.Atoi(rawBookID)
				if err == nil {
					bookID = &id
				}
			}

			var chapter *string
			rawChapter := r.FormValue("chapter")
			if rawChapter != "" {
				chapter = &rawChapter
			}

			postsModel := &models.PostModel{DB: h.DB}

			id, err := postsModel.InsertPost(userID, title, content, category, bookID, chapter)
			if err != nil {
				h.ErrorLog.Println("DB Error:", err)
				h.serverError(w, err)
				return
			}

			http.Redirect(w, r, fmt.Sprintf("/post/%d", id), http.StatusSeeOther)
		}
	}
}