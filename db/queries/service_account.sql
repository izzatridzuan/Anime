-- name: CreateServiceAccount :one
INSERT INTO service_accounts (name, api_key_hash)
VALUES ($1, $2)
RETURNING *;

-- name: ListServiceAccounts :many
SELECT * FROM service_accounts ORDER BY created_at DESC;

-- name: DeleteServiceAcount :exec
DELETE FROM service_accounts WHERE id = $1;

-- name: GetServiceAccountByKeyHash :one
SELECT * FROM service_accounts WHERE api_key_hash = $1;
