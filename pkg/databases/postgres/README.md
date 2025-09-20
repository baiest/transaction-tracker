# Database Management Guide with Migrations

This document explains how to manage changes to your PostgreSQL database schema using the `golang-migrate` tool. This guide assumes you have a running PostgreSQL instance that is accessible.

## 1\. Prerequisites

Make sure you have the following components installed:

  * **Go**: For the migration tool.
  * **`migrate` CLI**: The command-line tool for managing migrations.

You can install `migrate` with the following command:

```sh
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

-----

## 2\. Migration Management

Every change to the database schema should be a new migration. This ensures that all changes are versioned and can be rolled back.

### 2.1. Create a New Migration

To create a new pair of migration files (`.up.sql` and `.down.sql`), use the `migrate create` command. The name should be descriptive.

```sh
# Example for creating a new table
migrate create -ext sql -dir migrations -seq create_accounts_table

# Example for adding a new column
migrate create -ext sql -dir migrations -seq add_category_to_movements
```

### 2.2. Write the SQL Files

  * **`[...].up.sql`**: Contains the SQL statement to apply the change, such as `CREATE TABLE` or `ALTER TABLE`.
  * **`[...].down.sql`**: Contains the SQL statement to revert the change, such as `DROP TABLE` or `ALTER TABLE ... DROP COLUMN`.

### 2.3. Apply Migrations

To apply all pending migrations, use the `migrate up` command.

```sh
# To ensure the password isn't visible, use an environment variable
# On Linux/macOS: export PGPASSWORD="your_password"
# On Windows: set PGPASSWORD="your_password"

migrate -path migrations -database "postgres://user:password@localhost:5432/mydb?sslmode=disable" up
```

### 2.4. Roll Back Migrations

To revert the last applied migration, use `migrate down`.

```sh
migrate -path migrations -database "postgres://user:password@localhost:5432/mydb?sslmode=disable" down 1
```

-----

## 3\. Best Practices

  * **Environment Variables**: Always use environment variables for your database password (`PGPASSWORD`) to avoid hardcoding it in your commands or scripts.
  * **Data Integrity**: Use `FOREIGN KEY` to ensure referential integrity between tables.
  * **Optimization**: Create `INDEX` on columns that are frequently used for searching or filtering data.