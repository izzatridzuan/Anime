package handlers

import (
	"log/slog"
	"net/http"

	"github.com/unixadmin/anime/internal/db"
	audittmpl "github.com/unixadmin/anime/templates/audit"
)

type AuditLogHandler struct {
	queries *db.Queries
}

func NewAuditLogHandler(queries *db.Queries) *AuditLogHandler {
	return &AuditLogHandler{queries: queries}
}

func (h *AuditLogHandler) List(w http.ResponseWriter, r *http.Request) {
	logs, err := h.queries.ListAuditLogs(r.Context())
	if err != nil {
		slog.Error("failed to list audit logs", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	audittmpl.List(logs, getSessionRole(r)).Render(r.Context(), w)
}
