<div align="center">
  <img src="assets/logo.svg" alt="Postgresus Logo" width="250"/>
  
  <h3>PostgreSQL backup and restore tool</h3>
  <p>Free, open source and self-hosted solution for automated PostgreSQL backups with multiple storage options and notifications</p>
  
  <p>
    <a href="#-features">Features</a> •
    <a href="#-installation">Installation</a> •
    <a href="#-usage">Usage</a> •
    <a href="#-license">License</a> •
    <a href="#-contributing">Contributing</a>
  </p>
  
  <img src="assets/dashboard.svg" alt="Postgresus Dashboard" width="800"/>
</div>

---

## ✨ Features

### 🔄 **Scheduled Backups**

- **Flexible scheduling**: hourly, daily, weekly, monthly
- **Precise timing**: run backups at specific times (e.g., 4 AM during low traffic)
- **Smart compression**: 4-8x space savings with balanced compression (~20% overhead)

### 🗄️ **Multiple Storage Destinations**

- **Local storage**: Keep backups on your VPS/server
- **Cloud storage**: S3, Cloudflare R2, Google Drive, Dropbox, and more (coming soon)
- **Secure**: All data stays under your control

### 📱 **Smart Notifications**

- **Multiple channels**: Email, Telegram, Slack, webhooks (coming soon)
- **Real-time updates**: Success and failure notifications
- **Team integration**: Perfect for DevOps workflows

### 🐘 **PostgreSQL Support**

- **Multiple versions**: PostgreSQL 13, 14, 15, 16, and 17
- **SSL support**: Secure connections available
- **Easy restoration**: One-click restore from any backup

### 🐳 **Self-Hosted & Secure**

- **Docker-based**: Easy deployment and management
- **Privacy-first**: All your data stays on your infrastructure
- **Open source**: MIT licensed, inspect every line of code

---

## 📦 Installation

You have two ways to install Postgresus: via automated script (recommended) or manual Docker Compose setup.

### Option 1: Automated Installation Script (Recommended, Linux only)

The installation script will:

- ✅ Install Docker with Docker Compose (if not already installed)
- ✅ Create optimized `docker-compose.yml` configuration
- ✅ Set up automatic startup on system reboot via cron

```bash
sudo apt-get install -y curl && \
sudo curl -sSL https://raw.githubusercontent.com/RostislavDugin/postgresus/refs/heads/main/install-postgresus.sh \
| sudo bash
```

### Option 2: Manual Docker Compose Setup

Create a `docker-compose.yml` file with the following configuration:

```yaml
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
```

Then run:

```bash
docker compose up -d
```

---

## 🚀 Usage

1. **Access the dashboard**: Navigate to `http://localhost:4005`
2. **Add first DB for backup**: Click "New Database" and follow the setup wizard
3. **Configure schedule**: Choose from hourly, daily, weekly or monthly intervals
4. **Set database connection**: Enter your PostgreSQL credentials and connection details
5. **Choose storage**: Select where to store your backups (local, S3, Google Drive, etc.)
6. **Add notifications** (optional): Configure email, Telegram, Slack, or webhook notifications
7. **Save and start**: Postgresus will validate settings and begin the backup schedule

### 🔑 Resetting Admin Password

If you need to reset the admin password, you can use the built-in password reset command:

```bash
docker exec -it postgresus ./main --new-password="YourNewSecurePassword123"
```

---

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.
