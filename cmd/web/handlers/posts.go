package handlers

import (
    "database/sql"
    "errors"
    "net/http"
    "strconv"
    "html/template"
    "github.com/go-chi/chi/v5"
	"forum/internal/models"
	"fmt"
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
        IsAuthenticatedOk bool
    }{
        Post:              postView,
        IsAuthenticatedOk: h.isAuthenticated(r),
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
	postIDStr := chi.URLParam(r, "id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		return
	}

	err = r.ParseForm()
	if err != nil {
		h.clientError(w, http.StatusBadRequest)
		return
	}

	content := r.FormValue("content")
	if content == "" {
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