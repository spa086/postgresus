#!/bin/bash

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
      test: ["CMD-SHELL", "pg_isready -U postgresus -d postgresus -p 5437"]
      interval: 5s
      timeout: 5s
      retries: 5
    restart: unless-stopped
EOF
log "docker-compose.yml created successfully"

# Install cron if not installed
if ! command -v cron &> /dev/null; then
    log "Cron not found. Installing cron..."
    apt-get update
    apt-get install -y cron
    systemctl enable cron
    log "Cron installed successfully"
else
    log "Cron already installed"
fi

# Add cron job for system reboot
log "Setting up cron job to start PostgresUS on system reboot..."
CRON_JOB="@reboot cd $INSTALL_DIR && docker-compose up -d >> $INSTALL_DIR/postgresus-startup.log 2>&1"
(crontab -l 2>/dev/null | grep -v "$INSTALL_DIR.*docker-compose"; echo "$CRON_JOB") | crontab -
log "Cron job configured successfully"

log "PostgresUS installation completed successfully!"
log "-------------------------------------------"
log "To launch immediately:"
log "> cd $INSTALL_DIR && docker compose up -d"
log "Or reboot system to auto-start via cron"
log "Access PostgresUS at: http://localhost:4005"