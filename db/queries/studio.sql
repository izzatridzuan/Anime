-- name: GetStudio :one
SELECT * FROM studios WHERE id = $1 AND deleted_at IS NULL;

-- name: GetStudiosByAnime :many
SELECT s.* FROM studios s
JOIN anime_studios ans ON s.id = ans.studio_id
WHERE ans.anime_id = $1 
AND s.deleted_at IS NULL;

-- name: ListStudios :many
SELECT * FROM studios WHERE deleted_at IS NULL ORDER BY name ASC;

-- name: CreateStudio :one
INSERT INTO studios (name)
VALUES ($1)
RETURNING *;

-- name: ArchiveStudio :exec
UPDATE studios SET deleted_at = NOW() WHERE id = $1;

-- name: AddStudioToAnime :exec
INSERT INTO anime_studios (anime_id, studio_id)
VALUES ($1, $2);

-- name: RemoveStudioFromAnime :exec
DELETE FROM anime_studios
WHERE anime_id = $1 AND studio_id = $2;

-- name: DeleteAnimeStudios :exec
DELETE FROM anime_studios WHERE anime_id = $1;