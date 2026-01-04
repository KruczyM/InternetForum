package handlers

import (
	"fmt"
	"forum/internal/models"
	"forum/internal/validator"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
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
	userID := h.authenticatedUserID(r)
	if userID == "" {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}

	tab := r.URL.Query().Get("tab")
	if tab == "" {
		tab = "profile"
	}

	editMode := r.URL.Query().Get("edit")

	user, err := models.GetUserByID(h.DB, userID)
	if err != nil {
		h.serverError(w, err)
		return
	}

	data := h.newTemplateData(w,r)

	data.Form = &validator.Validator{}

	data.AnyData["tab"] = tab
	data.AnyData["editMode"] = editMode
	data.AnyData["user"] = user
	data.AnyData["postsCount"], _ = models.CountUserPosts(h.DB, userID)
	data.AnyData["commentsCount"], _ = models.CountUserComments(h.DB, userID)
	data.AnyData["likesCount"], _ = models.CountUserLikes(h.DB, userID)

	postsModel := &models.PostModel{DB: h.DB}

	switch tab {
	case "posts":
		data.AnyData["posts"], _ = postsModel.GetPostsByUserID(userID)

	case "comments":
		data.AnyData["comments"], _ = postsModel.GetCommentsByUserID(userID)

	case "likes":
		data.AnyData["likes"], _ = postsModel.GetLikesByUserID(userID)
	}

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
		data := h.newTemplateData(w,r)
		data.AnyData = map[string]any{}
		data.Form = form

		userID := h.authenticatedUserID(r)
		user, _ := models.GetUserByID(h.DB, userID)

		data.AnyData["tab"] = "profile"
		data.AnyData["editMode"] = "profile"
		data.AnyData["user"] = user

		h.render(w, http.StatusUnprocessableEntity, "user_panel.html", data)
		return
	}

	userID := h.authenticatedUserID(r)
	err = models.UpdateUserNameFields(h.DB, userID, form.firstName, form.lastName)
	if err != nil {
		h.serverError(w, err)
		return
	}
	h.setFlash(w, "success", "Profile updated successfully")

	http.Redirect(w, r, "/profile?tab=profile", http.StatusSeeOther)
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

	userID := h.authenticatedUserID(r)
	user, _ := models.GetUserByID(h.DB, userID)

	if !form.Valid() {
		data := h.newTemplateData(w,r)
		data.AnyData = map[string]any{}
		data.Form = form

		data.AnyData["tab"] = "profile"
		data.AnyData["editMode"] = "password"
		data.AnyData["user"] = user

		h.render(w, http.StatusUnprocessableEntity, "user_panel.html", data)
		return
	}

	if !models.CheckPassword(form.currentPassword, user.PasswordHash) {
		form.AddFieldError("current_password", "Incorrect current password")

		data := h.newTemplateData(w,r)
		data.AnyData = map[string]any{}
		data.Form = form

		data.AnyData["tab"] = "profile"
		data.AnyData["editMode"] = "password"
		data.AnyData["user"] = user

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
	h.setFlash(w,"success", "Password changed successfully")

	http.Redirect(w, r, "/profile?tab=profile", http.StatusSeeOther)
}

func (h *Handler) changeAvatar(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		h.serverError(w, err)
		return
	}

	var imagePath string

	file, handler, err := r.FormFile("avatar")
	if err != nil && err != http.ErrMissingFile {
		h.serverError(w, err)
		return
	}

	if file != nil {
		defer file.Close()

		fileName := fmt.Sprintf(
			"%d_%s",
			time.Now().UnixNano(),
			filepath.Base(handler.Filename),
		)

		uploadDir := "./ui/static/uploads/avatars"
		filePath := filepath.Join(uploadDir, fileName)

		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			h.serverError(w, err)
			return
		}

		dst, err := os.Create(filePath)
		if err != nil {
			h.serverError(w, err)
			return
		}
		defer dst.Close()

		if _, err := io.Copy(dst, file); err != nil {
			h.serverError(w, err)
			return
		}

		imagePath = "/static/uploads/avatars/" + fileName
	}

	userID := h.authenticatedUserID(r)

	err = models.UpdateUserAvatarPath(h.DB, userID, imagePath)
	if err != nil {
		h.serverError(w, err)
		return
	}
	h.setFlash(w, "success", "Avatart changed successfully")

	http.Redirect(w, r, "/profile?tab=profile", http.StatusSeeOther)
}

func (h *Handler) publicUserProfile(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/u/")
	parts := strings.Split(path, "/")

	if len(parts) != 1 || parts[0] == "" {
		h.notFound(w)
		return
	}

	username := parts[0]

    tab := r.URL.Query().Get("tab")
    if tab == "" {
        tab = "posts"
    }

    user, err := models.GetUserByUsername(h.DB, username)
    if err != nil {
        h.notFound(w)
        return
    }

    data := h.newTemplateData(w,r)
    data.AnyData["user"] = user
    data.AnyData["tab"] = tab

    postsModel := &models.PostModel{DB: h.DB}

    switch tab {
    case "posts":
        data.AnyData["posts"], _ = postsModel.GetPostsByUserID(user.ID)
    case "comments":
        data.AnyData["comments"], _ = postsModel.GetCommentsByUserID(user.ID)
    case "likes":
        data.AnyData["likes"], _ = postsModel.GetLikesByUserID(user.ID)
    }

    h.render(w, http.StatusOK, "public_user_panel.html", data)
}
