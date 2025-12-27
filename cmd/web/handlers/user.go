package handlers

import (
	"forum/internal/models"
	"forum/internal/validator"
	"net/http"
	"strings"
)

type userProfileForm struct {
	firstName string
	lastName  string
	validator.Validator
}

type userPasswordForm struct {
	currentPassword string
	newPassword     string
	confirmPassword string
	validator.Validator
}

func (h *Handler) userProfile(w http.ResponseWriter, r *http.Request) {
	userID := h.SessionManager.GetString(r.Context(), "authenticatedUserID")
	if userID == "" {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}

	user, err := models.GetUserByID(h.DB, userID)
	if err != nil {
		h.serverError(w, err)
		return
	}

	postsCount, _ := models.CountUserPosts(h.DB, userID)
	commentsCount, _ := models.CountUserComments(h.DB, userID)
	likesCount, _ := models.CountUserLikes(h.DB, userID)

	data := h.newTemplateData(r)
	data.Form = &validator.Validator{}
	data.AnyData["user"] = user
	data.AnyData["postsCount"] = postsCount
	data.AnyData["commentsCount"] = commentsCount
	data.AnyData["likesCount"] = likesCount
	data.AnyData["editMode"] = r.URL.Query().Get("edit")
	h.render(w, http.StatusOK, "user_panel.html", data)
}

func (h *Handler) userProfileEditPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		h.clientError(w, http.StatusBadRequest)
		return
	}

	form := &userProfileForm{
		firstName: strings.TrimSpace(r.PostForm.Get("first_name")),
		lastName:  strings.TrimSpace(r.PostForm.Get("last_name")),
	}

	form.CheckField(validator.NotBlank(form.firstName), "first_name", "First name is required")
	form.CheckField(validator.NotBlank(form.lastName), "last_name", "Last name is required")

	if !form.Valid() {
		data := h.newTemplateData(r)
		data.Form = form
		data.AnyData["editMode"] = "profile"
		h.render(w, http.StatusUnprocessableEntity, "user_panel.html", data)
		return
	}

	userID := h.SessionManager.GetString(r.Context(), "authenticatedUserID")
	err = models.UpdateUserNameFields(h.DB, userID, form.firstName, form.lastName)
	if err != nil {
		h.serverError(w, err)
		return
	}

	h.SessionManager.Put(r.Context(), "flash", &FlashMessage{
		Type: "success",
		Msg:  "Profile updated successfully",
	})

	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

func (h *Handler) userProfilePasswordPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		h.clientError(w, http.StatusBadRequest)
		return
	}

	form := &userPasswordForm{
		currentPassword: r.PostForm.Get("current_password"),
		newPassword:     r.PostForm.Get("new_password"),
		confirmPassword: r.PostForm.Get("confirm_password"),
	}

	form.CheckField(validator.NotBlank(form.currentPassword), "current_password", "Current password is required")
	form.CheckField(validator.NotBlank(form.newPassword), "new_password", "New password is required")
	form.CheckField(validator.MinChars(form.newPassword, 5), "new_password", "Password must be at least 5 characters")
	form.CheckField(form.newPassword == form.confirmPassword, "confirm_password", "Passwords do not match")

	if !form.Valid() {
		data := h.newTemplateData(r)
		data.Form = form
		data.AnyData["editMode"] = "password"
		h.render(w, http.StatusUnprocessableEntity, "user_panel.html", data)
		return
	}

	userID := h.SessionManager.GetString(r.Context(), "authenticatedUserID")
	user, err := models.GetUserByID(h.DB, userID)
	if err != nil {
		h.serverError(w, err)
		return
	}

	if !models.CheckPassword(form.currentPassword, user.PasswordHash) {
		form.AddFieldError("current_password", "Incorrect current password")

		data := h.newTemplateData(r)
		data.Form = form
		data.AnyData["editMode"] = "password"
		h.render(w, http.StatusUnprocessableEntity, "user_panel.html", data)
		return
	}

	hashed, err := models.HashPassword(form.newPassword)
	if err != nil {
		h.serverError(w, err)
		return
	}

	err = models.UpdateUserPasswordHash(h.DB, userID, hashed)
	if err != nil {
		h.serverError(w, err)
		return
	}

	h.SessionManager.Put(r.Context(), "flash", &FlashMessage{
		Type: "success",
		Msg:  "Password changed successfully",
	})

	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}
