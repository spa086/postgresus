
# Installation

You have 2 ways how to run Postgresus: install via script or manually write docker-compose.yml config.

**1) Install Postgresus via script (recommended, Linux only).**

It will:

- install Docker with Docker Compose
- write docker-compose.yml config
- install cron job to start Postgresus on system reboot

To install, run:

```
apt-get install -y curl && \
curl -sSL https://raw.githubusercontent.com/RostislavDugin/postgresus/refs/heads/main/install-postgresus.sh | bash
```

**2) Write docker-compose.yml config manually:**

```
version: "3"

services:
  postgresus:
    image: rostislavdugin/postgresus:latest
    ports:
      - "4005:4005"
    volumes:
      - ./postgresus-data:/app/postgresus-data
    depends_on:
      postgresus-db:
        condition: service_healthy
    restart: unless-stopped

  postgresus-db:
    image: postgres:17
    # we use default values, but do not expose
    # PostgreSQL ports so it is safe
    environment:
      - POSTGRES_DB=postgresus
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=Q1234567
    volumes:
      - ./pgdata:/var/lib/postgresql/data
    container_name: postgresus-db
    command: -p 5437
    shm_size: 10gb
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres -d postgresus -p 5437"]
      interval: 5s
      timeout: 5s
      retries: 5
    restart: unless-stopped
```

# Usage

Go to http://localhost:4005 to see Postgresus UI
