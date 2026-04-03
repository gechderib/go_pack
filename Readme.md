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

<!-- fo loggin use zap it what aws cloudwatch, elk and datadog epect for monitoring -->

go get go.uber.org/

<!-- advanced go -->
context.Context = request lifecycle controller

   It handles:
        cancellation
        timeout
        request-scoped values


Situation 	               What to do
DB query 	                    db.WithContext(ctx)
HTTP call	                    NewRequestWithContext(ctx)
Your own long task	          select { case <-ctx.Done() }
Goroutines / workers	     select { case <-ctx.Done() }


// i have loggin, recover, auth and timeout middlewares.
// in what order	should i use them?
// logging should be
// 1. first to log all requests,
// 2. then recovery to catch panics,
// 3. then auth to check authentication and
// 4. finally timeout to set a timeout for the request.
