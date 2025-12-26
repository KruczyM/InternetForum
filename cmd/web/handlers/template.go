package handlers

import (
	"forum/internal/models"
	"net/http"
)

type templateData struct {
	Post		*models.PostView
	Posts		[]models.PostView
	Books		[]models.Book
	CurrentYear int
	Flash		string
	IsAuthenticatedOk bool
	AuthenticatedUser string
	Form		any
}

func (h *Handler) newTemplateData(r *http.Request) *templateData {
	return &templateData {
		CurrentYear:	2025,
		IsAuthenticatedOk: h.isAuthenticated(r),
		AuthenticatedUser: h.SessionManager.GetString(r.Context(), "authenticatedUserID"),
	}
}