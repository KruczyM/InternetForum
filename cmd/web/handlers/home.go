package handlers

import (
	"fmt"
	"forum/internal/models"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func (h *Handler) home(w http.ResponseWriter, r *http.Request) {

	category := r.URL.Query().Get("category")
	bookIDStr := r.URL.Query().Get("book_id")
	query := r.URL.Query().Get("q")

	bookID := 0
	if bookIDStr != "" {
		if id, err := strconv.Atoi(bookIDStr); err == nil {
			bookID = id
		}
	}

	postsModel := &models.PostModel{DB: h.DB}

	var posts []models.PostView
	var err error

	if query != "" {
		posts, err = postsModel.SearchPosts(query, category)
	} else {
		posts, err = postsModel.GetAllPosts(category, bookID)
	}

	if err != nil {
		h.serverError(w, err)
		return
	}

	books, err := postsModel.GetAllBooks()
	if err != nil {
		h.serverError(w, err)
		return
	}

	data := struct {
		IsAuthenticated bool
		Posts           []models.PostView
		Books           []models.Book
		Query           string
	}{
		IsAuthenticated: h.isAuthenticated(r),
		Posts:           posts,
		Books:           books,
		Query:           query,
	}

	ts, err := template.ParseFiles("ui/html/home.html")
	if err != nil {
		h.serverError(w, err)
		return
	}
	ts.Execute(w, data)
}

func (h *Handler) createPost(w http.ResponseWriter, r *http.Request) {

	// ===== GET =====
	if r.Method == http.MethodGet {

		postsModel := &models.PostModel{DB: h.DB}
		books, err := postsModel.GetAllBooks()
		if err != nil {
			h.serverError(w, err)
			return
		}

		data := struct {
			Books           []models.Book
			IsAuthenticated bool
		}{
			Books:           books,
			IsAuthenticated: h.isAuthenticated(r),
		}

		ts, err := template.ParseFiles("ui/html/create.html")
		if err != nil {
			h.serverError(w, err)
			return
		}
		ts.Execute(w, data)
		return
	}

	// ===== POST =====
	if r.Method == http.MethodPost {

		if err := r.ParseMultipartForm(10 << 20); err != nil {
			h.serverError(w, err)
			return
		}

		title := r.FormValue("title")
		content := r.FormValue("content")
		category := r.FormValue("post_type")

		userID := h.SessionManager.GetString(r.Context(), "authenticatedUserID")
		if userID == "" {
			http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
			return
		}

		// --- book ---
		var bookID *int
		rawBookID := r.FormValue("book_id")
		if rawBookID != "" && rawBookID != "0" {
			if id, err := strconv.Atoi(rawBookID); err == nil {
				bookID = &id
			}
		}

		// --- chapter ---
		var chapter *string
		rawChapter := r.FormValue("chapter")
		if rawChapter != "" {
			chapter = &rawChapter
		}

		// --- image upload ---
		var imagePath string

		file, handler, err := r.FormFile("image")
		if err != nil && err != http.ErrMissingFile {
			h.serverError(w, err)
			return
		}

		if file != nil {
			defer file.Close()

			fileName := fmt.Sprintf(
				"%d_%s",
				time.Now().UnixNano(),
				filepath.Base(handler.Filename),
			)

			uploadDir := "./ui/static/uploads"
			filePath := filepath.Join(uploadDir, fileName)

			if err := os.MkdirAll(uploadDir, 0755); err != nil {
				h.serverError(w, err)
				return
			}

			dst, err := os.Create(filePath)
			if err != nil {
				h.serverError(w, err)
				return
			}
			defer dst.Close()

			if _, err := io.Copy(dst, file); err != nil {
				h.serverError(w, err)
				return
			}

			imagePath = "/static/uploads/" + fileName
		}

		postsModel := &models.PostModel{DB: h.DB}
		id, err := postsModel.InsertPost(
			userID,
			title,
			content,
			imagePath,
			category,
			bookID,
			chapter,
		)
		if err != nil {
			h.serverError(w, err)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/post/%d", id), http.StatusSeeOther)
	}
}
