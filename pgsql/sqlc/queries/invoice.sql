-- name: InvoiceHeaderCreate :one
INSERT INTO "invoice_header" (uuid, client_id) VALUES ($1, $2) RETURNING id, created_at;

-- name: InvoiceHeaderDeleteAll :exec
TRUNCATE TABLE "invoice_header" RESTART IDENTITY CASCADE;

-- name: InvoiceItemCreate :one
INSERT INTO "invoice_item" (invoice_header_id, product_id) VALUES ($1, $2) RETURNING id, created_at;

-- name: InvoiceItemDeleteAll :exec
TRUNCATE TABLE "invoice_item" RESTART IDENTITY;