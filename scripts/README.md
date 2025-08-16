# Database Scripts

This directory contains utility scripts for database operations.

## Scripts

### seed-dev-data.sql

Seeds the development database with sample user data for testing and development purposes.

**Usage:**

```bash
# Seed development database
psql -d graphql_service_dev -f scripts/seed-dev-data.sql

# Or using environment variables
PGDATABASE=graphql_service_dev psql -f scripts/seed-dev-data.sql
```

**Sample Data:**

- 10 sample users with realistic names and email addresses
- Varied creation and update timestamps to simulate real usage
- UUIDs that won't conflict with production data

## Environment Setup

Make sure your PostgreSQL database is running and accessible:

```bash
# Create development database (if not exists)
createdb graphql_service_dev

# Run migrations first
go run cmd/server/main.go -migrate

# Then seed data
psql -d graphql_service_dev -f scripts/seed-dev-data.sql
```

## Notes

- Seed data uses fixed UUIDs to ensure consistency across environments
- All operations use `ON CONFLICT DO NOTHING` to prevent duplicate key errors
- Scripts are idempotent and can be run multiple times safely
