package handlers

import (
	"forum/internal/models"
	"net/http"
)

type FlashMessage struct {
	Type string
	Msg  string
}

type templateData struct {
	Post		*models.PostView
	Posts		[]models.PostView
	Books		[]models.Book
	CurrentYear int
	Flash		*FlashMessage
	IsAuthenticatedOk bool
	AuthenticatedUser string
	Form		any
}

func (h *Handler) newTemplateData(r *http.Request) *templateData {
	data := &templateData {
		CurrentYear:	2025,
		IsAuthenticatedOk: h.isAuthenticated(r),
		AuthenticatedUser: h.SessionManager.GetString(r.Context(), "authenticatedUserID"),
	}

		flash := h.SessionManager.Pop(r.Context(), "flash")
	if f, ok := flash.(*FlashMessage); ok {
		data.Flash = f
	}
	return data
}

