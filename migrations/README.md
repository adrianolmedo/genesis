# Database Migrations

## How to run migrations

```bash
goose -dir ./migrations postgres "$DATABASE_URL" up
```

## How to rollback the last migration

```bash
goose -dir ./migrations postgres "$DATABASE_URL" down
```

## How to reset local dev database

```bash
psql "$DATABASE_URL" -f scripts/reset-db.sql
```
