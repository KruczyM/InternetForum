package handlers

import (
	"forum/internal/models"
	"html/template"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func (h *Handler) Register (w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		ts, err := template.ParseFiles("ui/html/register.html")
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal Server Error", 500)
			return
		}
		ts.Execute(w, nil)
		return
	}

	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		firstName := r.FormValue("first_name")
		lastName := r.FormValue("last_name")
		email := r.FormValue("email")
		password := r.FormValue("password")

		user := &models.User {
			ID: uuid.NewString(),
			Username: username,
			FirstName: firstName,
			LastName: lastName,
			Email: email,
			PasswordHash: models.HashPassword(password),
		}

		err := models.InsertUser(h.DB, user)
		if err != nil {
			log.Println("Registration failed:", err)
			http.Error(w, "Username or Email not valid", 400)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}