package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/internal/models"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

func (h *Handler) ViewPost(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/post/")
	parts := strings.Split(path, "/")

	if len(parts) != 1 {
		h.notFound(w, r)
		return
	}

	id, err := strconv.Atoi(parts[0])
	if err != nil || id < 1 {
		h.notFound(w, r)
		return
	}

	postsModel := &models.PostModel{DB: h.DB}

	postView, err := postsModel.GetPost(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.notFound(w, r)
		} else {
			h.serverError(w, err)
		}
		return
	}

	data := struct {
		Post              interface{}
		IsAuthenticated   bool
		AuthenticatedUser string
	}{
		Post:              postView,
		IsAuthenticated:   h.isAuthenticated(r),
		AuthenticatedUser: h.authenticatedUserID(r),
	}

	ts, err := template.ParseFiles("ui/html/post.html")
	if err != nil {
		h.ErrorLog.Println("Template Error:", err)
		h.serverError(w, err)
		return
	}
	ts.Execute(w, data)
}

func (h *Handler) CreateComment(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/post/")
	parts := strings.Split(path, "/")

	if len(parts) != 2 || parts[1] != "comment" {
		h.notFound(w, r)
		return
	}

	postID, err := strconv.Atoi(parts[0])
	if err != nil || postID < 1 {
		h.notFound(w, r)
		return
	}

	err = r.ParseForm()
	if err != nil {
		h.serverError(w, err)
		return
	}

	content := r.FormValue("content")
	if strings.TrimSpace(content) == "" {
		h.clientError(w, http.StatusBadRequest)
		return
	}

	userID := h.authenticatedUserID(r)
	if userID == "" {
		h.setFlash(w, "error", "Please log in to perform this action.")
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}

	postsModel := &models.PostModel{DB: h.DB}
	err = postsModel.InsertComment(postID, userID, content)
	if err != nil {
		h.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/post/%d", postID), http.StatusSeeOther)
}

func (h *Handler) DeletePost(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/post/")
	parts := strings.Split(path, "/")

	if len(parts) != 2 || parts[1] != "delete" {
		h.notFound(w, r)
		return
	}

	id, err := strconv.Atoi(parts[0])
	if err != nil || id < 1 {
		h.notFound(w, r)
		return
	}

	currentUserID := h.authenticatedUserID(r)
	postsModel := &models.PostModel{DB: h.DB}

	post, err := postsModel.GetPost(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.notFound(w, r)
		} else {
			h.serverError(w, err)
		}
		return
	}

	if post.UserID != currentUserID {
		h.clientError(w, http.StatusForbidden)
		return
	}

	err = postsModel.DeletePost(id)
	if err != nil {
		h.serverError(w, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) EditPost(w http.ResponseWriter, r *http.Request) {

	path := strings.TrimPrefix(r.URL.Path, "/post/")
	parts := strings.Split(path, "/")

	if len(parts) != 2 || parts[1] != "edit" {
		h.notFound(w, r)
		return
	}

	id, err := strconv.Atoi(parts[0])
	if err != nil || id < 1 {
		h.notFound(w, r)
		return
	}

	postModel := &models.PostModel{DB: h.DB}
	postView, err := postModel.GetPost(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.notFound(w, r)
		} else {
			h.serverError(w, err)
		}
		return
	}

	currentUserID := h.authenticatedUserID(r)
	if postView.Post.UserID != currentUserID {
		h.clientError(w, http.StatusForbidden)
		return
	}

	data := struct {
		Post            interface{}
		IsAuthenticated bool
	}{
		Post:            postView,
		IsAuthenticated: h.isAuthenticated(r),
	}

	ts, err := template.ParseFiles("ui/html/edit.html")
	if err != nil {
		h.serverError(w, err)
		return
	}
	ts.Execute(w, data)
}

func (h *Handler) UpdatePost(w http.ResponseWriter, r *http.Request) {

	path := strings.TrimPrefix(r.URL.Path, "/post/")
	parts := strings.Split(path, "/")

	if len(parts) != 2 || parts[1] != "edit" {
		h.notFound(w, r)
		return
	}

	id, err := strconv.Atoi(parts[0])
	if err != nil || id < 1 {
		h.notFound(w, r)
		return
	}

	err = r.ParseForm()
	if err != nil {
		h.serverError(w, err)
		return
	}

	title := r.Form.Get("title")
	content := r.Form.Get("content")

	postsModel := &models.PostModel{DB: h.DB}
	postView, err := postsModel.GetPost(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.notFound(w, r)
		} else {
			h.serverError(w, err)
		}
		return
	}

	currentUserID := h.authenticatedUserID(r)
	if postView.Post.UserID != currentUserID {
		h.clientError(w, http.StatusForbidden)
		return
	}

	err = postsModel.UpdatePost(id, title, content)
	if err != nil {
		h.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/post/%d", id), http.StatusSeeOther)
}

func (h *Handler) PostLike(w http.ResponseWriter, r *http.Request) {

	path := strings.TrimPrefix(r.URL.Path, "/post/")
	parts := strings.Split(path, "/")

	if len(parts) != 2 || parts[1] != "like" {
		h.notFound(w, r)
		return
	}

	id, err := strconv.Atoi(parts[0])
	if err != nil || id < 1 {
		h.notFound(w, r)
		return
	}

	userID := h.authenticatedUserID(r)

	if userID == "" {
		h.setFlash(w, "error", "Please log in to perform this action.")
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}

	postsModel := &models.PostModel{DB: h.DB}

	err = postsModel.ToggleLike(userID, id)
	if err != nil {
		h.serverError(w, err)
		return
	}

	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
}

// PostDislike toggles a dislike for a post, or switches from like to dislike.
func (h *Handler) PostDislike(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/post/")
	parts := strings.Split(path, "/")

	if len(parts) != 2 || parts[1] != "dislike" {
		h.notFound(w, r)
		return
	}

	id, err := strconv.Atoi(parts[0])
	if err != nil || id < 1 {
		h.notFound(w, r)
		return
	}

	userID := h.authenticatedUserID(r)
	if userID == "" {
		h.setFlash(w, "error", "Please log in to perform this action.")
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}

	postsModel := &models.PostModel{DB: h.DB}

	err = postsModel.ToggleDislike(userID, id)
	if err != nil {
		h.serverError(w, err)
		return
	}

	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
}

func (h *Handler) DeleteComment(w http.ResponseWriter, r *http.Request) {

	path := strings.TrimPrefix(r.URL.Path, "/comment/")
	parts := strings.Split(path, "/")

	if len(parts) != 2 || parts[1] != "delete" {
		h.notFound(w, r)
		return
	}

	id, err := strconv.Atoi(parts[0])
	if err != nil || id < 1 {
		h.notFound(w, r)
		return
	}

	userID := h.authenticatedUserID(r)
	if userID == "" {
		h.setFlash(w, "error", "Please log in to perform this action.")
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}

	postsModel := &models.PostModel{DB: h.DB}

	err = postsModel.DeleteComment(id, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			h.notFound(w, r)
		} else {
			h.serverError(w, err)
		}
		return
	}

	postID := r.URL.Query().Get("post_id")
	http.Redirect(w, r, "/post/"+postID, http.StatusSeeOther)
}

func (h *Handler) CommentLike(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/comment/")
	parts := strings.Split(path, "/")

	if len(parts) != 2 || parts[1] != "like" {
		h.notFound(w, r)
		return
	}

	id, err := strconv.Atoi(parts[0])
	if err != nil || id < 1 {
		h.notFound(w, r)
		return
	}

	userID := h.authenticatedUserID(r)
	if userID == "" {
		h.setFlash(w, "error", "Please log in to perform this action.")
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}

	postsModel := &models.PostModel{DB: h.DB}
	err = postsModel.ToggleLikeComment(userID, id)
	if err != nil {
		h.serverError(w, err)
		return
	}

	postID := r.URL.Query().Get("post_id")
	http.Redirect(w, r, "/post/"+postID, http.StatusSeeOther)
}

func (h *Handler) CommentDislike(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/comment/")
	parts := strings.Split(path, "/")

	if len(parts) != 2 || parts[1] != "dislike" {
		h.notFound(w, r)
		return
	}

	id, err := strconv.Atoi(parts[0])
	if err != nil || id < 1 {
		h.notFound(w, r)
		return
	}

	userID := h.authenticatedUserID(r)
	if userID == "" {
		h.setFlash(w, "error", "Please log in to perform this action.")
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}

	postsModel := &models.PostModel{DB: h.DB}
	err = postsModel.ToggleDislikeComment(userID, id)
	if err != nil {
		h.serverError(w, err)
		return
	}

	postID := r.URL.Query().Get("post_id")
	http.Redirect(w, r, "/post/"+postID, http.StatusSeeOther)
}

func (h *Handler) SearchPosts(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	category := r.URL.Query().Get("category")
	bookStr := r.URL.Query().Get("book")
	sort := r.URL.Query().Get("sort")

	bookID := 0
	if bookStr != "" {
		if id, err := strconv.Atoi(bookStr); err == nil {
			bookID = id
		}
	}

	postsModel := &models.PostModel{DB: h.DB}

	if strings.TrimSpace(query) == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	results, err := postsModel.SearchPosts(query, category, bookID, sort)
	if err != nil {
		h.serverError(w, err)
		return
	}

	data := struct {
		Posts             []models.PostView
		Query             string
		Category          string
		IsAuthenticated   bool
		AuthenticatedUser string
	}{
		Posts:             results,
		Query:             query,
		Category:          category,
		IsAuthenticated:   h.isAuthenticated(r),
		AuthenticatedUser: h.authenticatedUserID(r),
	}

	ts, err := template.ParseFiles("ui/html/search.html")
	if err != nil {
		h.serverError(w, err)
		return
	}

	ts.Execute(w, data)
}
