package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/unixadmin/anime/internal/db"
	"github.com/unixadmin/anime/templates"
	satmpl "github.com/unixadmin/anime/templates/service_account"
	"golang.org/x/crypto/bcrypt"
)

type ServiceAccountHandler struct {
	queries *db.Queries
}

func NewServiceAccountHandler(queries *db.Queries) *ServiceAccountHandler {
	return &ServiceAccountHandler{queries: queries}
}

func (h *ServiceAccountHandler) List(w http.ResponseWriter, r *http.Request) {
	accounts, err := h.queries.ListServiceAccounts(r.Context())
	if err != nil {
		slog.Error("failed to list service accounts", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	satmpl.List(accounts, getSessionRole(r), getFlash(w, r)).Render(r.Context(), w)
}

func (h *ServiceAccountHandler) New(w http.ResponseWriter, r *http.Request) {
	satmpl.Form(getSessionRole(r)).Render(r.Context(), w)
}

func (h *ServiceAccountHandler) Create(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	if name == "" {
		templates.AlertError("Name is required").Render(r.Context(), w)
		return
	}
	rawKey := generateAPIKey()
	hash, err := bcrypt.GenerateFromPassword([]byte(rawKey), bcrypt.DefaultCost)
	if err != nil {
		templates.AlertError("Failed to generate API key").Render(r.Context(), w)
		return
	}
	_, err = h.queries.CreateServiceAccount(r.Context(), db.CreateServiceAccountParams{
		Name:       name,
		ApiKeyHash: string(hash),
	})
	if err != nil {
		templates.AlertError("Name already exists").Render(r.Context(), w)
		return
	}

	satmpl.CreatedKey(rawKey, getSessionRole(r)).Render(r.Context(), w)
}

func (h *ServiceAccountHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))
	err := h.queries.DeleteServiceAcount(r.Context(), int32(id))
	if err != nil {
		slog.Error("failed to delete service account", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	setFlash(w, r, "Service account deleted successfully")
	w.Header().Set("HX-Redirect", "/service-accounts")
}

func generateAPIKey() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}
