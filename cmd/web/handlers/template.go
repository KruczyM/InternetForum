package handlers

import (
    "forum/internal/models"
    "net/http"
)
type templateData struct {
    Post              *models.PostView
    Posts             []models.PostView
    Books             []models.Book
    CurrentYear       int
    Flash             string
    IsAuthenticated   bool   
    AuthenticatedUser string
    Form              any
}

func (h *Handler) TemplateData(r *http.Request) *templateData {
    return &templateData {
        CurrentYear:       2025,
        IsAuthenticated:   h.isAuthenticated(r),
        AuthenticatedUser: h.SessionManager.GetString(r.Context(), "authenticatedUserID"),
        Flash:             h.SessionManager.PopString(r.Context(), "flash"),
    }
}