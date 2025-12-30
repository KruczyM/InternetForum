package handlers

import (
	"forum/internal/models"
	"html/template"
	"net/http"
	"strings"
)


func (h *Handler) CreateBook(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		data := h.newTemplateData(r)

		ts, err := template.ParseFiles("ui/html/create_book.html")
		if err != nil {
			h.serverError(w, err)
			return
		}
		ts.Execute(w, data)
		return
	}

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			h.clientError(w, http.StatusBadRequest)
			return
		}

		title := r.FormValue("title")
        author := r.FormValue("author")
        description := r.FormValue("description")

		if strings.TrimSpace(title) == "" || strings.TrimSpace(author) == "" {
			h.SessionManager.Put(r.Context(), "flash", "Title and Author are required!")
			http.Redirect(w, r, "/book/create", http.StatusSeeOther)
			return
		}

		bookModel := &models.BookModel{DB: h.DB}
		err = bookModel.AddBook(title, author, description)
		if err != nil {
			h.serverError(w, err)
			return
		}

		h.SessionManager.Put(r.Context(), "flash", "Book added successfully!")
		http.Redirect(w, r, "/post/create", http.StatusSeeOther)
	}
}