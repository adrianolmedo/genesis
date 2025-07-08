-- name: ProductCreate :one
INSERT INTO "product"
(uuid, name, observations, price, created_at)
VALUES ($1, $2, $3, $4, $5) RETURNING id;

-- name: ProductByID :one
SELECT * FROM "product" WHERE id = $1 AND deleted_at IS NULL;

-- name: ProductUpdate :one
UPDATE "product" 
SET 
    name = $1,
    observations = $2,
    price = $3,
    updated_at = $4
WHERE id = $5
RETURNING id;

-- name: ProductAll :many
SELECT * FROM "product" WHERE deleted_at IS NULL;

-- name: ProductDelete :one
UPDATE "product" SET deleted_at = $1 WHERE id = $2 RETURNING id;

-- name: ProductHardDelete :one
DELETE FROM "product" WHERE id = $1 RETURNING id;

-- name: ProductDeleteAll :exec
TRUNCATE TABLE "product" RESTART IDENTITY CASCADE;