package handlers

import (
	"fmt"
	"forum/internal/models"
	"forum/internal/validator"
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


type PostForm struct {
	Title    string
	Content  string
	Category string
	BookID   string
	Chapter  string
	Errors   map[string]string
}


func (h *Handler) createPost(w http.ResponseWriter, r *http.Request) {

	postsModel := &models.PostModel{DB: h.DB}

	if r.Method == http.MethodGet {

		books, err := postsModel.GetAllBooks()
		if err != nil {
			h.serverError(w, err)
			return
		}
		data := h.newTemplateData(w, r)
		data.AnyData["Books"] = books

		h.render(w, http.StatusOK, "create.html", data)
		return
	}

	if r.Method != http.MethodPost {
		h.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		h.clientError(w, http.StatusBadRequest)
		return
	}

	formTitle := r.FormValue("title")
	formContent := r.FormValue("content")
	formCategory := r.FormValue("post_type")
	formBookID := r.FormValue("book_id")
	formChapter := r.FormValue("chapter")

	v := validator.Validator{}


	v.CheckField(validator.NotBlank(formTitle),"title","Title cannot be blank",)
	v.CheckField(validator.MaxChars(formTitle, 150),"title","Title must be at most 150 characters",)
	v.CheckField(validator.NotBlank(formContent),"content","Content cannot be blank",)
	v.CheckField(validator.PermittedValue(
			formCategory,
			"Classics",
			"General_Fiction",
			"Crime_Mystery_Thriller",
			"Womens_Fiction",
			"Childrens",
			"Poetry_Plays",
			"Non_Fiction",
			"Historical_Fiction",
		),
		"post_type",
		"Invalid category selected",
	)

	var bookID *int
	if formBookID != "" && formBookID != "0" {
		id, err := strconv.Atoi(formBookID)
		if err != nil {
			v.AddFieldError("book_id", "Invalid book selected")
		} else {
			bookID = &id
		}
	}

	var chapter *string
	if validator.NotBlank(formChapter) {
		chapter = &formChapter
	}

	var imagePath string

	file, handler, err := r.FormFile("image")
	if err != nil && err != http.ErrMissingFile {
		v.AddFieldError("image", "Unable to upload image")
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

	if !v.Valid() {

		books, err := postsModel.GetAllBooks()
		if err != nil {
			h.serverError(w, err)
			return
		}

		data := h.newTemplateData(w, r)
		data.AnyData["Books"] = books
		data.AnyData["Form"] = map[string]any{
			"title":     formTitle,
			"content":   formContent,
			"post_type": formCategory,
			"book_id":   formBookID,
			"chapter":   formChapter,
			"errors":    v.FieldErrors,
		}

		w.WriteHeader(http.StatusBadRequest)
		h.render(w, http.StatusBadRequest, "create.html", data)
		return
	}

	userID := h.authenticatedUserID(r)
	if userID == "" {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}

	id, err := postsModel.InsertPost(
		userID,
		formTitle,
		formContent,
		imagePath,
		formCategory,
		bookID,
		chapter,
	)
	if err != nil {
		h.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/post/%d", id), http.StatusSeeOther)
}
