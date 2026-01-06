package handlers

import (
	"errors"
	"forum/internal/models"
	"forum/internal/validator"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type userRegisterForm struct {
	UserName  string
	FirstName string
	LastName  string
	Email     string
	Password  string
	validator.Validator
}

type userLoginForm struct {
	email    string
	password string
	validator.Validator
}

func (h *Handler) userRegister(w http.ResponseWriter, r *http.Request) {
	data := h.newTemplateData(w, r)
	data.Form = userRegisterForm{}
	h.render(w, http.StatusOK, "register.html", data)
}

func (h *Handler) userRegisterPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		h.clientError(w, http.StatusBadRequest)
		return
	}

	form := &userRegisterForm{
		Email:     strings.TrimSpace(r.PostForm.Get("email")),
		Password:  r.PostForm.Get("password"),
		FirstName: strings.TrimSpace(r.PostForm.Get("first_name")),
		LastName:  strings.TrimSpace(r.PostForm.Get("last_name")),
		UserName:  strings.TrimSpace(r.PostForm.Get("user_name")),
	}

	form.CheckField(validator.NotBlank(form.UserName), "user_name", "Username is required")
	form.CheckField(validator.NotBlank(form.FirstName), "first_name", "First name is required")
	form.CheckField(validator.NotBlank(form.LastName), "last_name", "Last name is required")
	form.CheckField(validator.NotBlank(form.Email), "email", "Email is required")
	form.CheckField(validator.NotBlank(form.Password), "password", "Password is required")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "Email is invalid")
	form.CheckField(validator.MinChars(form.Password, 5), "password", "Password must be at least 5 characters")

	if !form.Valid() {
		data := h.newTemplateData(w, r)
		data.Form = form
		h.render(w, http.StatusUnprocessableEntity, "register.html", data)
		return
	}

	hashedPassword, err := models.HashPassword(form.Password)
	if err != nil {
		h.ErrorLog.Println("Hashing password failed:", err)
		http.Error(w, "Internal Server Error", 500)
		return
	}

	user := &models.User{
		ID:           uuid.NewString(),
		Username:     form.UserName,
		FirstName:    form.FirstName,
		LastName:     form.LastName,
		Email:        form.Email,
		PasswordHash: hashedPassword,
	}

	err = models.InsertUser(h.DB, user)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrDuplicateEmail):
			form.AddFieldError("email", "Email address is already in use")

		case errors.Is(err, models.ErrDuplicateUsername):
			form.AddFieldError("user_name", "Username is already in use")

		default:
			h.serverError(w, err)
			return
		}

		data := h.newTemplateData(w, r)
		data.Form = form
		h.render(w, http.StatusUnprocessableEntity, "register.html", data)
		return
	}

	h.setFlash(w, "success", "Your signup was successful. Please log in.")

	http.Redirect(w, r, "/auth/login", http.StatusSeeOther)

}

func (h *Handler) userLogin(w http.ResponseWriter, r *http.Request) {
	data := h.newTemplateData(w, r)
	data.Form = userLoginForm{}
	h.render(w, http.StatusOK, "login.html", data)

}

func (h *Handler) userLoginPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		h.clientError(w, http.StatusBadRequest)
		return
	}

	form := &userLoginForm{
		email:    strings.TrimSpace(r.PostForm.Get("email")),
		password: r.PostForm.Get("password"),
	}

	form.CheckField(validator.NotBlank(form.email), "email", "Email is required")
	form.CheckField(validator.NotBlank(form.password), "password", "Password is required")

	if !form.Valid() {
		data := h.newTemplateData(w, r)
		data.Form = form
		h.render(w, http.StatusUnprocessableEntity, "login.html", data)
		return
	}

	user, err := models.GetUserByEmail(h.DB, form.email)
	if err != nil {
		h.setFlash(w, "error", "Invalid email or password")
		h.ErrorLog.Println(err)
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}

	if !models.CheckPassword(form.password, user.PasswordHash) {
		h.ErrorLog.Println("Invalid password")
		h.setFlash(w, "error", "Invalid email or password")
		h.ErrorLog.Println(err)
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}

	h.setUserSession(w, user.ID)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) userLogoutPost(w http.ResponseWriter, r *http.Request) {

	h.clearUserSession(w)
	h.setFlash(w, "success", "You've been logged out successfully!")

	http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
}
