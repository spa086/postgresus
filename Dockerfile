# ========= BUILD FRONTEND =========
FROM --platform=$BUILDPLATFORM node:24-alpine AS frontend-build

WORKDIR /frontend

COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci
COPY frontend/ ./

# Copy .env file (with fallback to .env.production.example)
RUN if [ ! -f .env ]; then \
      if [ -f .env.production.example ]; then \
        cp .env.production.example .env; \
      fi; \
    fi

RUN npm run build

# ========= BUILD BACKEND =========
FROM --platform=$BUILDPLATFORM golang:1.23.3 AS backend-build

# Install Go public tools needed in runtime
RUN curl -fsSL https://raw.githubusercontent.com/pressly/goose/master/install.sh | sh
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Set working directory
WORKDIR /app

# Install Go dependencies
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Create required directories for embedding
RUN mkdir -p /app/ui/build

# Copy frontend build output for embedding
COPY --from=frontend-build /frontend/dist /app/ui/build

# Generate Swagger documentation
COPY backend/ ./
RUN swag init -d . -g cmd/main.go -o swagger

# Compile the backend
ARG TARGETOS
ARG TARGETARCH
ARG TARGETVARIANT
RUN CGO_ENABLED=0 \
    GOOS=$TARGETOS \
    GOARCH=$TARGETARCH \
    go build -o /app/main ./cmd/main.go


# ========= RUNTIME =========
FROM --platform=$TARGETPLATFORM debian:bookworm-slim

# Install PostgreSQL client tools (versions 13-17)
RUN apt-get update && apt-get install -y --no-install-recommends \
       wget ca-certificates gnupg lsb-release && \
    wget -qO- https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add - && \
    echo "deb http://apt.postgresql.org/pub/repos/apt $(lsb_release -cs)-pgdg main" \
      > /etc/apt/sources.list.d/pgdg.list && \
    apt-get update && \
    apt-get install -y --no-install-recommends \
       postgresql-client-13 postgresql-client-14 postgresql-client-15 \
       postgresql-client-16 postgresql-client-17 && \
    rm -rf /var/lib/apt/lists/*

# Create symlinks for PostgreSQL client tools
RUN for v in 13 14 15 16 17; do \
      mkdir -p /usr/pgsql-$v/bin && \
      for b in pg_dump psql pg_restore createdb dropdb; do \
        ln -sf /usr/bin/$b /usr/pgsql-$v/bin/$b; \
      done; \
    done

WORKDIR /app

# Copy Goose from build stage
COPY --from=backend-build /usr/local/bin/goose /usr/local/bin/goose

# Copy app binary 
COPY --from=backend-build /app/main .

# Copy migrations directory
COPY backend/migrations ./migrations

# Copy UI files
COPY --from=backend-build /app/ui/build ./ui/build

# Copy .env file (with fallback to .env.production.example)
COPY backend/.env* /app/
RUN if [ ! -f /app/.env ]; then \
      if [ -f /app/.env.production.example ]; then \
        cp /app/.env.production.example /app/.env; \
      fi; \
    fi

EXPOSE 4005

CMD ["./main"]
