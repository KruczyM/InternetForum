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
	"strings"
	"time"
)

func (h *Handler) home(w http.ResponseWriter, r *http.Request) {

	// --- query params ---
	query := strings.TrimSpace(r.URL.Query().Get("q"))
	category := r.URL.Query().Get("category")
	bookStr := r.URL.Query().Get("book")
	sort := r.URL.Query().Get("sort")

	bookID := 0
	if bookStr != "" {
		if id, err := strconv.Atoi(bookStr); err == nil {
			bookID = id
		}
	}

	postModel := &models.PostModel{DB: h.DB}
	bookModel := &models.BookModel{DB: h.DB}

	posts, err := postModel.SearchPosts(query, category, bookID, sort)
	if err != nil {
		h.serverError(w, err)
		return
	}

	books, err := bookModel.GetAllBooks()
	if err != nil {
		h.serverError(w, err)
		return
	}

	data := struct {
		IsAuthenticated bool
		Posts           []models.PostView
		Books           []models.Book
		Query           string
		Category        string
		SelectedBook    int
		Sort            string
	}{
		IsAuthenticated: h.isAuthenticated(r),
		Posts:           posts,
		Books:           books,
		Query:           query,
		Category:        category,
		SelectedBook:    bookID,
		Sort:            r.URL.Query().Get("sort"),
	}

	ts, err := template.ParseFiles("ui/html/home.html")
	if err != nil {
		h.serverError(w, err)
		return
	}

	err = ts.Execute(w, data)
	if err != nil {
		h.serverError(w, err)
	}
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

		userID := h.authenticatedUserID(r)
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

			uploadDir := "./ui/static/uploads/posts_img"
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

			imagePath = "/static/uploads/posts_img/" + fileName
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
