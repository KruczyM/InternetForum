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
	userName  string
	firstName string
	lastName  string
	email     string
	password  string
	validator.Validator
}

type userLoginForm struct {
	email    string
	password string
	validator.Validator
}

// i make i sepatare to make it easier to read
func (h *Handler) userRegister(w http.ResponseWriter, r *http.Request) {
	data := h.newTemplateData(r)
	data.Form = userRegisterForm{}
	h.render(w, http.StatusOK, "register.html", data)
}

func (h *Handler) userRegisterPost(w http.ResponseWriter, r *http.Request) {
	//more secure is to use http.Request.ParseForm and then r.PostForm.Get("") instead of r.FormValue("")
	err := r.ParseForm()
	if err != nil {
		h.clientError(w, http.StatusBadRequest)
		return
	}

	form := &userRegisterForm{
		email:     strings.TrimSpace(r.PostForm.Get("email")),
		password:  r.PostForm.Get("password"),
		firstName: strings.TrimSpace(r.PostForm.Get("first_name")),
		lastName:  strings.TrimSpace(r.PostForm.Get("last_name")),
		userName:  strings.TrimSpace(r.PostForm.Get("user_name")),
	}

	form.CheckField(validator.NotBlank(form.userName), "user_name", "Username is required")
	form.CheckField(validator.NotBlank(form.firstName), "first_name", "First name is required")
	form.CheckField(validator.NotBlank(form.lastName), "last_name", "Last name is required")
	form.CheckField(validator.NotBlank(form.email), "email", "Email is required")
	form.CheckField(validator.NotBlank(form.password), "password", "Password is required")
	form.CheckField(validator.Matches(form.email, validator.EmailRX), "email", "Email is invalid")
	form.CheckField(validator.MinChars(form.password, 5), "password", "Password must be at least 5 characters")

	if !form.Valid() {
		data := h.newTemplateData(r)
		data.Form = form
		h.render(w, http.StatusUnprocessableEntity, "register.html", data)
		return
	}

	hashedPassword, err := models.HashPassword(form.password)
	if err != nil {
		h.ErrorLog.Println("Hashing password failed:", err)
		http.Error(w, "Internal Server Error", 500)
		return
	}

	user := &models.User{
		ID:           uuid.NewString(),
		Username:     form.userName,
		FirstName:    form.firstName,
		LastName:     form.lastName,
		Email:        form.email,
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

		data := h.newTemplateData(r)
		data.Form = form
		h.render(w, http.StatusUnprocessableEntity, "register.html", data)
		return
	}

	h.SessionManager.Put(r.Context(), "flash", &FlashMessage{
		Type: "success",
		Msg:  "Your signup was successful. Please log in.",
	})
	http.Redirect(w, r, "/auth/login", http.StatusSeeOther)

}

func (h *Handler) userLogin(w http.ResponseWriter, r *http.Request) {
	data := h.newTemplateData(r)
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
		data := h.newTemplateData(r)
		data.Form = form
		h.render(w, http.StatusUnprocessableEntity, "login.html", data)
		return
	}

	user, err := models.GetUserByEmail(h.DB, form.email)
	if err != nil {
		h.SessionManager.Put(r.Context(), "flash", &FlashMessage{
			Type: "error",
			Msg:  "Invalid email or password",
		})
		h.ErrorLog.Println(err)
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}

	if !models.CheckPassword(form.password, user.PasswordHash) {
		h.ErrorLog.Println("Invalid password")
		h.SessionManager.Put(r.Context(), "flash", &FlashMessage{
			Type: "error",
			Msg:  "Invalid email or password",
		})
		h.ErrorLog.Println(err)
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
	h.SessionManager.Put(r.Context(), "flash", &FlashMessage{
		Type: "success",
		Msg:  "You've been logged out successfully!",
	})
	// Redirect the user to the application home page.
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
