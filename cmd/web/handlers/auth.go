package handlers

import (
	"forum/internal/models"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

// i make i sepatare to make it easier to read
func (h *Handler) userRegister(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("ui/html/register.html")
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", 500)
		return
	}
	ts.Execute(w, nil)
	return
}

func (h *Handler) userRegisterPost(w http.ResponseWriter, r *http.Request) {

	username := r.FormValue("username")
	firstName := r.FormValue("first_name")
	lastName := r.FormValue("last_name")
	email := r.FormValue("email")
	password := r.FormValue("password")

	hashedPassword, err := models.HashPassword(password)
	if err != nil {
		h.ErrorLog.Println("Hashing password failed:", err)
		http.Error(w, "Internal Server Error", 500)
		return
	}

	user := &models.User{
		ID:           uuid.NewString(),
		Username:     username,
		FirstName:    firstName,
		LastName:     lastName,
		Email:        email,
		PasswordHash: hashedPassword,
	}
	err = models.InsertUser(h.DB, user)
	if err != nil {
		h.ErrorLog.Println("Registration failed:", err)
		http.Error(w, "Username or Email not valid", 400)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func (h *Handler) userLogin(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("ui/html/login.html")
	if err != nil {
		h.ErrorLog.Println("Template Error:", err)
		h.serverError(w, err)
		return
	}

	flash := h.SessionManager.PopString(r.Context(), "flash")
	data := struct{ Flash string }{Flash: flash}
	if err := tmpl.Execute(w, data); err != nil {
		h.serverError(w, err)
	}
}

func (h *Handler) userLoginPost(w http.ResponseWriter, r *http.Request) {
	email := strings.TrimSpace(r.FormValue("email"))
	password := r.FormValue("password")

	user, err := models.GetUserByEmail(h.DB, email)
	if err != nil {
		h.SessionManager.Put(r.Context(), "flash", "Invalid email or password")
		h.ErrorLog.Println(err)
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}

	if !models.CheckPassword(password, user.PasswordHash) {
		h.ErrorLog.Println("Invalid password")
		h.SessionManager.Put(r.Context(), "flash", "Invalid email or password")
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}

	h.SessionManager.Put(r.Context(), "authenticatedUserID", user.ID)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	// Use the RenewToken() method on the current session to change the session
	// ID again.
	err := h.SessionManager.RenewToken(r.Context())
	if err != nil {
		h.serverError(w, err)
		return
	}
	// Remove the authenticatedUserID from the session data so that the user is
	// 'logged out'.
	h.SessionManager.Remove(r.Context(), "authenticatedUserID")
	// Add a flash message to the session to confirm to the user that they've been
	// logged out.
	h.SessionManager.Put(r.Context(), "flash", "You've been logged out successfully!")
	// Redirect the user to the application home page.
	http.Redirect(w, r, "/", http.StatusSeeOther)
}