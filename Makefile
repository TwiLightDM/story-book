USER := "postgres_user"
PASSWORD := "superStrongPassword"
HOST := "localhost"
PORT := "5432"
DB_NAME := "story-book"

DB_DSN := "postgres://${USER}:${PASSWORD}@${HOST}:${PORT}/${DB_NAME}?sslmode=disable"
MIGRATE := migrate -path ./migrations -database $(DB_DSN)

migrate-new:
	migrate create -ext sql -dir ./migrations ${NAME}
migrate:
	$(MIGRATE) up
migrate-down:
	$(MIGRATE) down
migrate-force:
	$(MIGRATE) force ${VERSION}

run:
	go run main.go

swagger:
	swag init -g main.go -o internal/docs

lint:
	golangci-lint run