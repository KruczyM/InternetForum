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
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("ui/html/register.html")
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", 500)
		return
	}
	ts.Execute(w, nil)
	return
}

func (h *Handler) RegisterPost(w http.ResponseWriter, r *http.Request) {

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

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("ui/html/login.html")
	if err != nil {
		h.ErrorLog.Println("Template Error:", err)
		h.serverError(w, err)
		return
	}
	template.Execute(w, nil)
}

func (h *Handler) LoginPost(w http.ResponseWriter, r *http.Request) {
	email := strings.TrimSpace(r.FormValue("email"))
	password := r.FormValue("password")

	user, err := models.GetUserByEmail(h.DB, email)
	if err != nil {
		h.ErrorLog.Println("DB Error:", err)
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
	}

	if !models.CheckPassword(password, user.PasswordHash) {
		h.clientError(w, http.StatusUnauthorized)
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
	}

	//h.SessionManager.Put(r.Context(), "userID", user.ID)
	http.Redirect(w, r, "/", http.StatusSeeOther)

}
