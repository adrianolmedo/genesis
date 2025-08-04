#!/bin/bash

ENV_FILE=".env.development"
RESET_SQL="scripts/reset-db.sql"

# Verify that the .env.development file exists.
if [[ ! -f "$ENV_FILE" ]]; then
  echo "❌ ERROR: $ENV_FILE not found"
  exit 1
fi

# Load environment variables from the ENV_FILE.
set -o allexport
source "$ENV_FILE"
set +o allexport

# Strict configuration errors.
set -euo pipefail

# Only allow reset in non-production environments.
ENV=${ENV:-development}

if [[ "$ENV" == "production" ]]; then
  echo "❌ ERROR: Reset is not allowed in production."
  exit 1
fi

# Verify that goose is installed.
if ! command -v goose &> /dev/null; then
  echo "❌ ERROR: goose not found in PATH."
  exit 1
fi

# Verify that exists the reset SQL file.
if [[ ! -f "$RESET_SQL" ]]; then
  echo "❌ ERROR: $RESET_SQL not found."
  exit 1
fi

# Manual confirmation before proceeding.
read -p "⚠️  This will DROP ALL DATA. Are you sure? (y/N): " confirm
confirm=${confirm,,}  # convertir a minúsculas
if [[ "$confirm" != "y" && "$confirm" != "yes" ]]; then
  echo "Cancelled."
  exit 0
fi

echo "🔄 Resetting database..."
echo "Executing $RESET_SQL"

sudo -u postgres psql -d "$DBNAME" -f "$RESET_SQL"

# Replicate migrations.
goose "$DBENGINE" "$DBURL" up

echo "✅ Database reset complete."
