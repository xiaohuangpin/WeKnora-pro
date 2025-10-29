#!/bin/bash
set -e

# Database connection details (can be overridden by environment variables)
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_USER=${DB_USER:-postgres}
DB_PASSWORD=${DB_PASSWORD:-postgres}
DB_NAME=${DB_NAME:-WeKnora}
MIGRATIONS_DIR="/app/migrations"

# Check if migrate tool is installed
if ! command -v migrate &> /dev/null; then
    echo "Error: migrate tool is not installed"
    echo "Install it with: go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest"
    exit 1
fi

# Construct the database URL
DB_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"

# Execute migration based on command
case "$1" in
    up)
        echo "Running migrations up..."
        echo "DB_URL: ${DB_URL}"
        echo "DB_USER: ${DB_USER}"
        echo "DB_PASSWORD: ${DB_PASSWORD}"
        echo "DB_HOST: ${DB_HOST}"
        echo "DB_PORT: ${DB_PORT}"
        echo "DB_NAME: ${DB_NAME}"
        echo "MIGRATIONS_DIR: ${MIGRATIONS_DIR}"
        migrate -path ${MIGRATIONS_DIR} -database ${DB_URL} up
        ;;
    down)
        echo "Running migrations down..."
        migrate -path ${MIGRATIONS_DIR} -database ${DB_URL} down
        ;;
    create)
        if [ -z "$2" ]; then
            echo "Error: Migration name is required"
            echo "Usage: $0 create <migration_name>"
            exit 1
        fi
        echo "Creating migration files for $2..."
        migrate create -ext sql -dir ${MIGRATIONS_DIR} -seq $2
        ;;
    *)
        echo "Usage: $0 {up|down|create <migration_name>}"
        exit 1
        ;;
esac

echo "Migration command completed successfully" 