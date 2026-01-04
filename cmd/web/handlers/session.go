package handlers

import (
	"net/http"
	"strings"
)

func (h *Handler) setUserSession(w http.ResponseWriter, userID string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_user",
		Value:    userID,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

func (h *Handler) clearUserSession(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:   "session_user",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
}

func (h *Handler) getUserID(r *http.Request) string {
	c, err := r.Cookie("session_user")
	if err != nil {
		return ""
	}
	return c.Value
}

func (h *Handler) setFlash(w http.ResponseWriter, flashType, msg string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "flash",
		Value:    flashType + "|" + msg,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   5,
	})
}

func (h *Handler) getFlash(w http.ResponseWriter, r *http.Request) *FlashMessage {
	c, err := r.Cookie("flash")
	if err != nil {
		return nil
	}

	parts := strings.SplitN(c.Value, "|", 2)
	if len(parts) != 2 {
		return nil
	}

	http.SetCookie(w, &http.Cookie{
		Name:   "flash",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	return &FlashMessage{
		Type: parts[0],
		Msg:  parts[1],
	}
}

func (h *Handler) authenticatedUserID(r *http.Request) string {
	id, _ := r.Context().Value("authenticatedUserID").(string)
	return id
}