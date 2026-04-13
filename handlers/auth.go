package handlers

import (
	"net/http"

	"github.com/unixadmin/anime/internal/db"
	"github.com/unixadmin/anime/templates/auth"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	queries *db.Queries
}

func NewAuthHandler(queries *db.Queries) *AuthHandler {
	return &AuthHandler{queries: queries}
}

func (h *AuthHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	auth.Login().Render(r.Context(), w)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	user, err := h.queries.GetUserByEmail(r.Context(), email)
	if err != nil {
		auth.LoginError("Invalid email or password").Render(r.Context(), w)
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		auth.LoginError("Invalid email or password").Render(r.Context(), w)
		return
	}

	session, _ := store.Get(r, "session")
	session.Values["user_id"] = int(user.ID)
	session.Values["role"] = user.Role
	session.Save(r, w)

	if user.MustChangePassword {
		w.Header().Set("HX-Redirect", "/change-password")
		return
	}
	w.Header().Set("HX-Redirect", "/anime")
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	session.Options.MaxAge = -1
	session.Save(r, w)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (h *AuthHandler) ChangePasswordPage(w http.ResponseWriter, r *http.Request) {
	auth.ChangePassword().Render(r.Context(), w)
}

func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	userID := getSessionUserID(r)
	password := r.FormValue("password")
	confirm := r.FormValue("confirm_password")

	if err := validatePasswordComplexity(password); err != nil {
		auth.ChangePasswordError(err.Message).Render(r.Context(), w)
		return
	}

	if password != confirm {
		auth.ChangePasswordError("Passwords do not match").Render(r.Context(), w)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		auth.ChangePasswordError("Failed to update password").Render(r.Context(), w)
		return
	}

	err = h.queries.UpdateUserPassword(r.Context(), db.UpdateUserPasswordParams{
		ID:           int32(userID),
		PasswordHash: string(hash),
	})
	if err != nil {
		auth.ChangePasswordError("Failed to update password").Render(r.Context(), w)
		return
	}

	w.Header().Set("HX-Redirect", "/anime")
}
