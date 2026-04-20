package handlers

import (
	"log/slog"
	"net/http"

	"github.com/unixadmin/anime/internal/db"
	"github.com/unixadmin/anime/templates"
	profiletmpl "github.com/unixadmin/anime/templates/profile"
	"golang.org/x/crypto/bcrypt"
)

type ProfileHandler struct {
	queries *db.Queries
}

func NewProfileHandler(queries *db.Queries) *ProfileHandler {
	return &ProfileHandler{queries: queries}
}

func (h *ProfileHandler) Page(w http.ResponseWriter, r *http.Request) {
	userID := getSessionUserID(r)
	user, err := h.queries.GetUserByID(r.Context(), int32(userID))
	if err != nil {
		slog.Error("user not found on profile page", "error", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	profiletmpl.Profile(user, getSessionRole(r)).Render(r.Context(), w)
}

func (h *ProfileHandler) UpdateName(w http.ResponseWriter, r *http.Request) {
	userID := getSessionUserID(r)
	name := r.FormValue("name")

	if name == "" {
		templates.AlertError("Name cannot be empty").Render(r.Context(), w)
		return
	}

	err := h.queries.UpdateUserName(r.Context(), db.UpdateUserNameParams{
		ID:   int32(userID),
		Name: name,
	})
	if err != nil {
		templates.AlertError("Failed to update name").Render(r.Context(), w)
		return
	}
	templates.AlertSuccess("Name updated successfully").Render(r.Context(), w)
}

func (h *ProfileHandler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	userID := getSessionUserID(r)
	currentPassword := r.FormValue("current_password")
	newPassword := r.FormValue("password")
	confirm := r.FormValue("confirm_password")

	user, err := h.queries.GetUserByID(r.Context(), int32(userID))
	if err != nil {
		templates.AlertError("User not found").Render(r.Context(), w)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(currentPassword)); err != nil {
		templates.AlertError("Current password is incorrect").Render(r.Context(), w)
		return
	}

	if err := validatePasswordComplexity(newPassword); err != nil {
		templates.AlertError(err.Message).Render(r.Context(), w)
		return
	}

	if newPassword != confirm {
		templates.AlertError("Passwords do not match").Render(r.Context(), w)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		templates.AlertError("Failed to update password").Render(r.Context(), w)
		return
	}

	err = h.queries.UpdateUserPassword(r.Context(), db.UpdateUserPasswordParams{
		ID:           int32(userID),
		PasswordHash: string(hash),
	})

	if err != nil {
		templates.AlertError("Failed to update password").Render(r.Context(), w)
		return
	}
	templates.AlertSuccess("Password changed successfully").Render(r.Context(), w)
}
