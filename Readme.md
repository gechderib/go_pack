## About migrations(manual migration id db schema changes):

go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
migrate -version

# 👉 Key idea:
up.sql → apply change
down.sql → rollback change

# Generate migration
migrate create -ext sql -dir migrations -seq create_users_table
migrate create -ext sql -dir migrations -seq create_orders_table
        will Create 
             ....up.sql
             ....down.sql

# Edit up.sql and down.sql which contain your current change it can be create, alert, drop ...


# Run migrations
SQLite
migrate -path migrations -database "sqlite3://test.db" up

Postgres
migrate -path migrations \
-database "postgres://myuser:mypass@localhost:5436/mydb?sslmode=disable" \
up

# Rollback
migrate -path migrations -database "..." down 1

👉 Rolls back 1 step

<!-- Learn the Following go orm -->
sqlx or Bun > GORM