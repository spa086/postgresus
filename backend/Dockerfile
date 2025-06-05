FROM golang:1.23.3

EXPOSE 4005

WORKDIR /app

# Install PostgreSQL APT repository and all versions 13-17
RUN apt-get update && apt-get install -y wget ca-certificates gnupg lsb-release
RUN wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add -
RUN echo "deb http://apt.postgresql.org/pub/repos/apt/ $(lsb_release -cs)-pgdg main" > /etc/apt/sources.list.d/pgdg.list
RUN apt-get update

# Install PostgreSQL client tools for versions 13-17
RUN apt-get install -y \
    postgresql-client-13 \
    postgresql-client-14 \
    postgresql-client-15 \
    postgresql-client-16 \
    postgresql-client-17

# Create symlinks to match expected paths (/usr/pgsql-{VERSION}/bin)
RUN for version in 13 14 15 16 17; do \
        mkdir -p /usr/pgsql-${version}/bin; \
        ln -sf /usr/bin/pg_dump /usr/pgsql-${version}/bin/pg_dump; \
        ln -sf /usr/bin/psql /usr/pgsql-${version}/bin/psql; \
        ln -sf /usr/bin/pg_restore /usr/pgsql-${version}/bin/pg_restore; \
        ln -sf /usr/bin/createdb /usr/pgsql-${version}/bin/createdb; \
        ln -sf /usr/bin/dropdb /usr/pgsql-${version}/bin/dropdb; \
    done

# Install global dependencies
RUN curl -fsSL https://raw.githubusercontent.com/pressly/goose/master/install.sh | sh
RUN goose -version
RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN swag init -d . -g cmd/main.go -o swagger
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/main.go

CMD ["./main"]