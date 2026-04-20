package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/unixadmin/anime/email"
	"github.com/unixadmin/anime/internal/db"
	"github.com/unixadmin/anime/templates"
	usertmpl "github.com/unixadmin/anime/templates/user"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	queries *db.Queries
}

func NewUserHandler(queries *db.Queries) *UserHandler {
	return &UserHandler{queries: queries}
}

func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	var users []db.User
	var err error
	if search != "" {
		users, err = h.queries.SearchUsers(r.Context(), search)
	} else {
		users, err = h.queries.ListUsers(r.Context())
	}
	if err != nil {
		slog.Error("failed to list users", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	usertmpl.List(users, search, getSessionRole(r), getFlash(w, r)).Render(r.Context(), w)
}

func (h *UserHandler) New(w http.ResponseWriter, r *http.Request) {
	usertmpl.Form(db.User{}, false, getSessionRole(r)).Render(r.Context(), w)
}

func (h *UserHandler) Edit(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))
	user, err := h.queries.GetUserByID(r.Context(), int32(id))
	if err != nil {
		slog.Error("user not found on edit", "error", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	usertmpl.Form(user, true, getSessionRole(r)).Render(r.Context(), w)
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))
	name := r.FormValue("name")
	userEmail := r.FormValue("email")
	role := r.FormValue("role")

	if err := validateUserForm(name, userEmail, role); err != nil {
		templates.AlertError(err.Message).Render(r.Context(), w)
		return
	}

	_, err := h.queries.UpdateUser(r.Context(), db.UpdateUserParams{
		ID:    int32(id),
		Name:  name,
		Email: userEmail,
		Role:  role,
	})
	if err != nil {
		templates.AlertError("Email already exists").Render(r.Context(), w)
		return
	}
	logAudit(r.Context(), h.queries, getSessionUserID(r), "updated", "user", int32(id))
	setFlash(w, r, "User updated successfully")
	w.Header().Set("HX-Redirect", "/users")
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	userEmail := r.FormValue("email")
	role := r.FormValue("role")

	if err := validateUserForm(name, userEmail, role); err != nil {
		templates.AlertError(err.Message).Render(r.Context(), w)
		return
	}

	password := generatePassword()
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		templates.AlertError("Failed to create user").Render(r.Context(), w)
		return
	}

	newUser, err := h.queries.CreateUser(r.Context(), db.CreateUserParams{
		Email:        userEmail,
		Name:         name,
		PasswordHash: string(hash),
		Role:         role,
	})
	if err != nil {
		templates.AlertError("Email already exists").Render(r.Context(), w)
		return
	}

	if err := email.SendWelcomeEmail(userEmail, name, password); err != nil {
		templates.AlertError("User created but failed to send email").Render(r.Context(), w)
		return
	}
	logAudit(r.Context(), h.queries, getSessionUserID(r), "created", "user", newUser.ID)
	setFlash(w, r, "User created successfully")
	w.Header().Set("HX-Redirect", "/users")
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))
	err := h.queries.DeleteUser(r.Context(), int32(id))
	if err != nil {
		slog.Error("failed to delete user", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logAudit(r.Context(), h.queries, getSessionUserID(r), "archived", "user", int32(id))
	setFlash(w, r, "User deleted successfully")
	w.Header().Set("HX-Redirect", "/users")
}

func (h *UserHandler) ResendInvite(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))
	user, err := h.queries.GetUserByID(r.Context(), int32(id))
	if err != nil {
		slog.Error("user not found on resend", "error", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	password := generatePassword()
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("failed to reset password", "error", err)
		http.Error(w, "Failed to reset password", http.StatusInternalServerError)
		return
	}
	err = h.queries.ResetUserPassword(r.Context(), db.ResetUserPasswordParams{
		ID:           user.ID,
		PasswordHash: string(hash),
	})
	if err != nil {
		slog.Error("failed to send invite email", "error", err)
		http.Error(w, "Failed to reset password", http.StatusInternalServerError)
		return
	}
	if err := email.SendWelcomeEmail(user.Email, user.Name, password); err != nil {
		templates.AlertError("Failed to send email").Render(r.Context(), w)
		return
	}
	templates.AlertSuccess("Invite email sent successfully").Render(r.Context(), w)
}

func generatePassword() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)[:16] + "A1!"
}
