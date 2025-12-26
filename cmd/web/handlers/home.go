package handlers

import (
	"fmt"
	"forum/internal/models"
	"html/template"
	"net/http"
	"strconv"
)

func (h *Handler) home(w http.ResponseWriter, r *http.Request) {

	category := r.URL.Query().Get("category")
	bookIDStr := r.URL.Query().Get("book_id")
	bookID := 0

	if bookIDStr != "" {
		id, err := strconv.Atoi(bookIDStr)
		if err == nil {
			bookID = id
		}
	}

	postsModel := &models.PostModel{DB: h.DB}

	posts, err := postsModel.GetAllPosts(category, bookID)
	if err != nil {
		h.ErrorLog.Println("DB Error:", err)
		h.serverError(w, err)
		return
	}

	books, err := postsModel.GetAllBooks()
	if err != nil {
		h.ErrorLog.Println("DB Error:", err)
		h.serverError(w, err)
		return
	}

	data := struct {
		IsAuthenticatedOk	bool
		Posts []models.PostView
		Books []models.Book
	}{
		Posts: posts,
		Books: books,
		IsAuthenticatedOk: h.isAuthenticated(r),
	}

	ts, err := template.ParseFiles("ui/html/home.html")
	if err != nil {
		h.ErrorLog.Println("Template Error:", err)
		h.serverError(w, err)
		return
	}
	ts.Execute(w, data)
}

/*func (h *Handler) createPost(w http.ResponseWriter, r *http.Request) {
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

			err := r.ParseForm()
			if err != nil {
				h.serverError(w, err)
				return
			}

			title := r.FormValue("title")
			content := r.FormValue("content")
			category := r.FormValue("category")
			userID := h.SessionManager.GetString(r.Context(), "authenticatedUserID")

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

			fmt.Printf("DEBUG: UserID is '%s\n", userID)

			postsModel := &models.PostModel{DB: h.DB}

			id, err := postsModel.InsertPost(userID, title, content, category, bookID, chapter)
			if err != nil {
				h.ErrorLog.Println("DB Error:", err)
				h.serverError(w, err)
				return
			}

			http.Redirect(w, r, fmt.Sprintf("/post/%d", id), http.StatusSeeOther)
	}
}*/

func (h *Handler) createPost(w http.ResponseWriter, r *http.Request) {

    if r.Method == http.MethodGet {
        
        postsModel := &models.PostModel{DB: h.DB}
        books, err := postsModel.GetAllBooks()
        if err != nil {
            h.serverError(w, err)
            return
        }

        data := struct {
            Books []models.Book
            IsAuthenticatedOk bool
        }{
            Books: books,
            IsAuthenticatedOk: h.isAuthenticated(r),
        }

        ts, err := template.ParseFiles("ui/html/create.html")
        if err != nil {
            h.ErrorLog.Println("Template Error:", err)
            h.serverError(w, err)
            return
        }
        ts.Execute(w, data)
        return
    }

    if r.Method == http.MethodPost {

        err := r.ParseForm()
        if err != nil {
            h.serverError(w, err)
            return
        }

        title := r.FormValue("title")
        content := r.FormValue("content")
        
        category := r.FormValue("post_type") 
        
        userID := h.SessionManager.GetString(r.Context(), "authenticatedUserID")

        var bookID *int
        rawBookID := r.FormValue("book_id")

        if rawBookID != "" && rawBookID != "0" {
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