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
	Post              *models.PostView
	Posts             []models.PostView
	Books             []models.Book
	CurrentYear       int
	Flash             *FlashMessage
	IsAuthenticated   bool
	AuthenticatedUser string
	Form              any
	AnyData           map[string]any
}

func (h *Handler) newTemplateData(w http.ResponseWriter, r *http.Request) *templateData {
	data := &templateData{
		CurrentYear:       2025,
		IsAuthenticated: h.isAuthenticated(r),
		AuthenticatedUser: h.authenticatedUserID(r),
		AnyData:           make(map[string]any),
		Flash: h.getFlash(w, r),
	}

	return data
}
