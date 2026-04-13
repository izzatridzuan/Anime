-- name: GetAnime :one
SELECT * FROM anime WHERE id = $1 AND deleted_at IS NULL;

-- name: ListAnime :many
SELECT * FROM anime WHERE  deleted_at IS NULL ORDER BY created_at DESC;

-- name: CreateAnime :one
INSERT INTO anime (title, genre, episodes, status, image_url, release_date)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: UpdateAnime :one
UPDATE anime
SET title = $2, genre = $3, episodes = $4, status = $5, image_url = $6, release_date = $7
WHERE id = $1
RETURNING *;


-- name: ArchiveAnime :exec
UPDATE anime SET deleted_at = NOW() WHERE id = $1;

-- name: FilterAnimePaginated :many
SELECT * FROM anime
WHERE
    ($1::text = '' OR LOWER(title) LIKE LOWER('%' || $1 || '%'))
    AND ($2::text = '' OR LOWER(genre) LIKE LOWER('%' || $2 || '%'))
    AND ($3::text = '' OR status = $3)
    AND ($4::timestamp IS NULL OR release_date >= $4)
    AND ($5::timestamp IS NULL OR release_date <= $5)
    AND deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $6 OFFSET $7;

-- name: CountFilteredAnime :one
SELECT COUNT(*) FROM anime
WHERE
    ($1::text = '' OR LOWER(title) LIKE LOWER('%' || $1 || '%'))
    AND ($2::text = '' OR LOWER(genre) LIKE LOWER('%' || $2 || '%'))
    AND ($3::text = '' OR status = $3)
    AND ($4::timestamp IS NULL OR release_date >= $4)
    AND ($5::timestamp IS NULL OR release_date <= $5)
    AND deleted_at IS NULL;