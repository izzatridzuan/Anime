package handlers

import (
	"context"
	"log"

	"github.com/unixadmin/anime/internal/db"
)

func logAudit(ctx context.Context, queries *db.Queries, userID int, actionType, entityType string, entityID int32) {
	err := queries.CreateAuditLog(ctx, db.CreateAuditLogParams{
		UserID:     int32(userID),
		ActionType: actionType,
		EntityType: entityType,
		EntityID:   entityID,
	})
	if err != nil {
		log.Printf("audit log failed: %v", err)
	}
}
