package handlers

import (
	"fmt"
	"forum/internal/models"
	"html/template"
	"net/http"
	"strconv"
	"io"
	"path/filepath"
	"time"
	"os"
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

func (h *Handler) createPost(w http.ResponseWriter, r *http.Request) {

    if r.Method == http.MethodGet {
        postsModel := &models.PostModel{DB: h.DB}
        books, err := postsModel.GetAllBooks()
        if err != nil {
            h.serverError(w, err)
            return
        }

        data := struct {
            Books             []models.Book
            IsAuthenticated   bool
        }{
            Books:           books,
            IsAuthenticated: h.isAuthenticated(r),
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

        err := r.ParseMultipartForm(10 << 20)
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

        var imagePath string

        file, handler, err := r.FormFile("image")
        if err != nil && err != http.ErrMissingFile {
            h.serverError(w, err)
            return
        }

        if file != nil {
            defer file.Close()

            fileName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), filepath.Base(handler.Filename))
            
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
        
        id, err := postsModel.InsertPost(userID, title, content, imagePath, category, bookID, chapter)
        if err != nil {
            h.ErrorLog.Println("DB Error:", err)
            h.serverError(w, err)
            return
        }

        http.Redirect(w, r, fmt.Sprintf("/post/%d", id), http.StatusSeeOther)
    }
}