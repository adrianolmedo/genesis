-- name: UserCreate :one
INSERT INTO "user" 
(uuid, first_name, last_name, email, password, created_at)
VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;

-- name: UserByLogin :one
SELECT id FROM "user" WHERE email = $1 AND password = $2 AND deleted_at IS NULL;

-- name: UserByID :one
SELECT * FROM "user" WHERE id = $1 AND deleted_at IS NULL;

-- name: UserUpdate :exec
UPDATE "user" SET first_name = $1, last_name = $2, email = $3, password = $4, updated_at = $5 WHERE id = $6;

-- name: UserAll :many
SELECT * FROM "user" WHERE deleted_at IS NULL;

-- name: UserDelete :exec
UPDATE "user" SET deleted_at = $1 WHERE id = $2;

-- name: UserHardDelete :exec
DELETE FROM "user" WHERE id = $1;

-- name: UserDeleteAll :exec
TRUNCATE TABLE "user" RESTART IDENTITY;