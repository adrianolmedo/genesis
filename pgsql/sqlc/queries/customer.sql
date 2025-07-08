-- name: CustomerCreate :one
INSERT INTO "customer" (
    uuid,
    first_name,
    last_name,
    email,
    password,
    created_at
) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;

-- name: CustomerListAsc :many
SELECT * FROM "customer" WHERE deleted_at IS NULL ORDER BY @sort::text ASC LIMIT $1 OFFSET $2;

-- name: CustomerListCount :one
SELECT COUNT (*) FROM "customer" WHERE deleted_at IS NULL;

-- name: CustomerDelete :one
UPDATE "customer" SET deleted_at = $1 WHERE id = $2 RETURNING id;