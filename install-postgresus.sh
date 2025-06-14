#!/bin/bash

# Check if script is run as root
if [ "$(id -u)" -ne 0 ]; then
    echo "Error: This script must be run as root (sudo ./install-postgresus.sh)" >&2
    exit 1
fi

# Set up logging
LOG_FILE="/var/log/postgresus-install.log"
INSTALL_DIR="/opt/postgresus"

log() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] $1" | tee -a "$LOG_FILE"
}

# Create log file if doesn't exist
touch "$LOG_FILE"
log "Starting PostgresUS installation..."

# Create installation directory
log "Creating installation directory..."
if [ ! -d "$INSTALL_DIR" ]; then
    mkdir -p "$INSTALL_DIR"
    log "Created directory: $INSTALL_DIR"
else
    log "Directory already exists: $INSTALL_DIR"
fi

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    log "Docker not found. Installing Docker..."
    
    # Install Docker
    apt-get update
    apt-get remove -y docker docker-engine docker.io containerd runc
    apt-get install -y ca-certificates curl gnupg lsb-release
    mkdir -p /etc/apt/keyrings
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
    echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
    apt-get update
    apt-get install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin
    log "Docker installed successfully"
else
    log "Docker already installed"
fi

# Check if docker-compose is installed
if ! command -v docker-compose &> /dev/null && ! command -v docker compose &> /dev/null; then
    log "Docker Compose not found. Installing Docker Compose..."
    apt-get update
    apt-get install -y docker-compose-plugin
    log "Docker Compose installed successfully"
else
    log "Docker Compose already installed"
fi

# Write docker-compose.yml
log "Writing docker-compose.yml to $INSTALL_DIR"
cat > "$INSTALL_DIR/docker-compose.yml" << 'EOF'
version: "3"

services:
  postgresus:
    container_name: postgresus
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
    container_name: postgresus-db
    image: postgres:17
    # we use default values, but do not expose
    # PostgreSQL ports so it is safe
    environment:
      - POSTGRES_DB=postgresus
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=Q1234567
    volumes:
      - ./pgdata:/var/lib/postgresql/data
    command: -p 5437
    shm_size: 10gb
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres -d postgresus -p 5437"]
      interval: 5s
      timeout: 5s
      retries: 5
    restart: unless-stopped
EOF
log "docker-compose.yml created successfully"

# Start PostgresUS
log "Starting PostgresUS..."
cd "$INSTALL_DIR"
docker compose up -d
log "PostgresUS started successfully"

log "Postgresus installation completed successfully!"
log "-------------------------------------------"
log "To launch:"
log "> cd $INSTALL_DIR && docker compose up -d"
log "Access Postgresus at: http://localhost:4005"