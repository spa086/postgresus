# Before run

Keep in mind: you need to use dev-db from docker-compose.yml in this folder
instead of postgresus-db from docker-compose.yml in the root folder.

> Copy .env.example to .env
> Copy docker-compose.yml.example to docker-compose.yml (for development only)

# Run

To build:

> go build /cmd/main.go

To run:

> go run /cmd/main.go

To run tests:

> go test ./internal/... 

Before commit (make sure `golangci-lint` is installed):

> golangci-lint fmt
> golangci-lint run

# Migrations

To create migration:

> goose create MIGRATION_NAME sql

To run migrations:

> goose up

If latest migration failed:

To rollback on migration:

> goose down

# Swagger

To generate swagger docs:

> swag init -g .\cmd\main.go -o swagger

Swagger URL is:

> http://localhost:8080/api/v1/docs/swagger/index.html#/

# Project structure

Default endpoint structure is:

/feature
/feature/controller.go
/feature/service.go
/feature/repository.go
/feature/model.go
/feature/dto.go

If there are couple of models:
/feature/models/model1.go
/feature/models/model2.go
...

# Project rules

Always use time.Now().UTC() instead of time.Now()
