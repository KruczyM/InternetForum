package handlers

import (
	"forum/internal/models"
	"net/http"
)

func (h *Handler) homepage(w http.ResponseWriter, r *http.Request) {
	categoryModel := &models.CategoryModel{DB: h.DB}

	categories, err := categoryModel.GetCategoryOverview()
	if err != nil {
		h.serverError(w, err)
		return
	}

	data := h.newTemplateData(w,r)
	data.AnyData["categories"] = categories

	h.render(w, http.StatusOK, "homepage.html", data)
}
