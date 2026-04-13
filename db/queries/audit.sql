-- name: CreateAuditLog :exec
INSERT INTO audit_logs (user_id, action_type, entity_type, entity_id)
VALUES ($1, $2, $3, $4);

-- name: ListAuditLogs :many
SELECT al.id, al.user_id, u.name AS user_name, al.action_type, al.entity_type, al.entity_id, al.created_at
FROM audit_logs al
JOIN users u ON al.user_id = u.id
ORDER BY al.created_at DESC
LIMIT 100;
