package handlers

import (
	"net/http"
	"strconv"

	"github.com/unixadmin/anime/internal/db"
	"github.com/unixadmin/anime/templates"
	studiotmpl "github.com/unixadmin/anime/templates/studio"
)

type StudioHandler struct {
	queries *db.Queries
}

func NewStudioHandler(queries *db.Queries) *StudioHandler {
	return &StudioHandler{queries: queries}
}

func (h *StudioHandler) List(w http.ResponseWriter, r *http.Request) {
	studios, err := h.queries.ListStudios(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	studiotmpl.List(studios, getSessionRole(r), getFlash(w, r)).Render(r.Context(), w)
}

func (h *StudioHandler) New(w http.ResponseWriter, r *http.Request) {
	studiotmpl.Form(getSessionRole(r)).Render(r.Context(), w)
}

func (h *StudioHandler) Create(w http.ResponseWriter, r *http.Request) {
	if err := validateStudioForm(r.FormValue("name")); err != nil {
		templates.AlertError(err.Message).Render(r.Context(), w)
		return
	}
	studio, err := h.queries.CreateStudio(r.Context(), r.FormValue("name"))
	if err != nil {
		templates.AlertError("Failed to create studio. Please try again.").Render(r.Context(), w)
		return
	}
	logAudit(r.Context(), h.queries, getSessionUserID(r), "created", "studio", studio.ID)
	setFlash(w, r, "Studio created successfully")
	w.Header().Set("HX-Redirect", "/studios")
}

func (h *StudioHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))
	err := h.queries.ArchiveStudio(r.Context(), int32(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logAudit(r.Context(), h.queries, getSessionUserID(r), "archived", "studio", int32(id))
	setFlash(w, r, "Studio deleted successfully")
	w.Header().Set("HX-Redirect", "/studios")
}
