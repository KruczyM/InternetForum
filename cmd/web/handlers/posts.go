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

	"github.com/go-chi/chi/v5"
)

func (h *Handler) ViewPost(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		h.notFound(w)
		return
	}

	postsModel := &models.PostModel{DB: h.DB}

	postView, err := postsModel.GetPost(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.notFound(w)
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
		AuthenticatedUser: h.SessionManager.GetString(r.Context(), "authenticatedUserID"),
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
	idStr := chi.URLParam(r, "id")
	postID, err := strconv.Atoi(idStr)
	if err != nil || postID < 1 {
		h.notFound(w)
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

	userID := h.SessionManager.GetString(r.Context(), "authenticatedUserID")
	if userID == "" {
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
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		h.notFound(w)
		return
	}

	currentUserID := h.SessionManager.GetString(r.Context(), "authenticatedUserID")

	postsModel := &models.PostModel{DB: h.DB}

	post, err := postsModel.GetPost(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.notFound(w)
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
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		h.notFound(w)
		return
	}

	postModel := &models.PostModel{DB: h.DB}
	postView, err := postModel.GetPost(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.notFound(w)
		} else {
			h.serverError(w, err)
		}
		return
	}

	currentUserID := h.SessionManager.GetString(r.Context(), "authenticatedUserID")
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
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		h.notFound(w)
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
			h.notFound(w)
		} else {
			h.serverError(w, err)
		}
		return
	}

	currentUserID := h.SessionManager.GetString(r.Context(), "authenticatedUserID")
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

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || id < 1 {
		h.notFound(w)
		return
	}

	userID := h.SessionManager.GetString(r.Context(), "authenticatedUserID")
	if userID == "" {
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

func (h *Handler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	fmt.Println("--- DEBUG: HANDLER WAS HIT! ---")
	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		h.notFound(w)
		return
	}

	userID := h.SessionManager.GetString(r.Context(), "authenticatedUserID")
	if userID == "" {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}

	postsModel := &models.PostModel{DB: h.DB}

	err = postsModel.DeleteComment(id, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			h.notFound(w)
		} else {
			h.serverError(w, err)
		}
		return
	}

	postID := r.URL.Query().Get("post_id")
	http.Redirect(w, r, "/post/"+postID, http.StatusSeeOther)
}

func (h *Handler) CommentLike(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		h.notFound(w)
		return
	}

	userID := h.SessionManager.GetString(r.Context(), "authenticatedUserID")
	if userID == "" {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}

	postsModel := &models.PostModel{DB: h.DB}
	err = postsModel.LikeComment(id, userID)
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

	postsModel := &models.PostModel{DB: h.DB}

	if strings.TrimSpace(query) == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	results, err := postsModel.SearchPosts(query, category)
	if err != nil {
		h.serverError(w, err)
		return
	}

	data := struct {
		Posts             []models.PostView
		Query             string
		IsAuthenticated   bool
		AuthenticatedUser string
	}{
		Posts:             results,
		Query:             query,
		IsAuthenticated:   h.isAuthenticated(r),
		AuthenticatedUser: h.SessionManager.GetString(r.Context(), "authenticatedUserID"),
	}

	ts, err := template.ParseFiles("ui/html/search.html")
	if err != nil {
		h.serverError(w, err)
		return
	}

	ts.Execute(w, data)
}
