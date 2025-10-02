# Deployment Guide - Ubuntu Server

This guide covers deploying the voidling Discord bot on an Ubuntu server with automatic restart on crash.

## Prerequisites

- Ubuntu 20.04 LTS or newer
- `sudo` access
- Go 1.24.1 or higher (optional, if building on server)

## Option 1: Deploy Pre-built Binary (Recommended)

### Step 1: Build the Binary Locally

On your development machine:

```bash
# Build for Linux
make build-linux

# The binary will be in build/voidling
```

### Step 2: Transfer to Server

```bash
# Using scp
scp build/voidling user@your-server:/home/user/voidling/

# Or using rsync
rsync -avz build/voidling user@your-server:/home/user/voidling/
```

### Step 3: Set Up on Server

SSH into your server and set up the directory structure:

```bash
ssh user@your-server

# Create directory structure
mkdir -p ~/voidling
cd ~/voidling

# Make binary executable
chmod +x voidling

# Create .env file
nano .env
```

Add your configuration to `.env`:

```bash
DISCORD_TOKEN=your_bot_token_here
DATABASE_PATH=/home/user/voidling/data/voidling.db
LOG_LEVEL=info
```

```bash
# Create data directory
mkdir -p data

# Test run
./voidling
```

If it starts successfully, press `Ctrl+C` and proceed to set up systemd.

## Option 2: Build on Server

### Step 1: Install Go on Ubuntu

```bash
# Download Go (check for latest version)
wget https://go.dev/dl/go1.24.1.linux-amd64.tar.gz

# Remove old Go installation (if exists)
sudo rm -rf /usr/local/go

# Extract new Go
sudo tar -C /usr/local -xzf go1.24.1.linux-amd64.tar.gz

# Add to PATH (add to ~/.bashrc for persistence)
export PATH=$PATH:/usr/local/go/bin
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc

# Verify installation
go version
```

### Step 2: Clone and Build

```bash
# Clone repository
cd ~
git clone https://github.com/kaffeed/voidling.git
cd voidling

# Install dependencies
go mod download

# Build
go build -trimpath -ldflags "-s -w" -o voidling ./cmd/voidling

# Create .env file
cp .env.example .env
nano .env
```

Add your Discord token to `.env`:

```bash
DISCORD_TOKEN=your_bot_token_here
DATABASE_PATH=/home/user/voidling/data/voidling.db
LOG_LEVEL=info
```

```bash
# Create data directory
mkdir -p data

# Test run
./voidling
```

## Set Up systemd Service (Automatic Restart)

This will make the bot start automatically on server boot and restart if it crashes.

### Step 1: Create systemd Service File

```bash
sudo nano /etc/systemd/system/voidling.service
```

Add the following content (adjust paths for your username):

```ini
[Unit]
Description=Voidling Discord Bot
After=network.target

[Service]
Type=simple
User=your-username
WorkingDirectory=/home/your-username/voidling
ExecStart=/home/your-username/voidling/voidling
Restart=always
RestartSec=10
StandardOutput=append:/home/your-username/voidling/logs/voidling.log
StandardError=append:/home/your-username/voidling/logs/voidling-error.log

# Environment file
EnvironmentFile=/home/your-username/voidling/.env

# Security hardening
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=read-only
ReadWritePaths=/home/your-username/voidling/data /home/your-username/voidling/logs

[Install]
WantedBy=multi-user.target
```

**Important:** Replace `your-username` with your actual Linux username.

### Step 2: Create Log Directory

```bash
mkdir -p ~/voidling/logs
```

### Step 3: Reload systemd and Enable Service

```bash
# Reload systemd to recognize new service
sudo systemctl daemon-reload

# Enable service to start on boot
sudo systemctl enable voidling

# Start the service
sudo systemctl start voidling

# Check status
sudo systemctl status voidling
```

## Managing the Service

### Check Status

```bash
sudo systemctl status voidling
```

### Start/Stop/Restart

```bash
# Start
sudo systemctl start voidling

# Stop
sudo systemctl stop voidling

# Restart
sudo systemctl restart voidling
```

### View Logs

```bash
# Real-time logs
sudo journalctl -u voidling -f

# Last 100 lines
sudo journalctl -u voidling -n 100

# Logs since today
sudo journalctl -u voidling --since today

# Application logs (from files)
tail -f ~/voidling/logs/voidling.log
tail -f ~/voidling/logs/voidling-error.log
```

### Disable Auto-start

```bash
sudo systemctl disable voidling
```

## Updating the Bot

### Method 1: Replace Binary (Pre-built)

```bash
# On local machine, build new version
make build-linux

# Transfer to server
scp build/voidling user@your-server:/home/user/voidling/voidling-new

# On server
ssh user@your-server
cd ~/voidling

# Stop service
sudo systemctl stop voidling

# Backup old binary
mv voidling voidling.backup

# Replace with new binary
mv voidling-new voidling
chmod +x voidling

# Restart service
sudo systemctl start voidling

# Check status
sudo systemctl status voidling
```

### Method 2: Git Pull and Rebuild (Built on Server)

```bash
cd ~/voidling

# Stop service
sudo systemctl stop voidling

# Pull latest code
git pull origin main

# Rebuild
go build -trimpath -ldflags "-s -w" -o voidling ./cmd/voidling

# Run database migrations (happens automatically on start)

# Restart service
sudo systemctl start voidling

# Check status
sudo systemctl status voidling
```

## Database Backups

### Manual Backup

```bash
# Create backup directory
mkdir -p ~/voidling/backups

# Backup database
cp ~/voidling/data/voidling.db ~/voidling/backups/voidling-$(date +%Y%m%d-%H%M%S).db
```

### Automated Daily Backups (Cron)

```bash
# Edit crontab
crontab -e
```

Add this line to backup daily at 2 AM:

```bash
0 2 * * * cp /home/your-username/voidling/data/voidling.db /home/your-username/voidling/backups/voidling-$(date +\%Y\%m\%d).db
```

Keep only last 7 days of backups:

```bash
0 3 * * * find /home/your-username/voidling/backups -name "voidling-*.db" -mtime +7 -delete
```

## Monitoring

### Check if Bot is Running

```bash
# Check systemd service
sudo systemctl is-active voidling

# Check process
ps aux | grep voidling

# Check listening ports (if any)
sudo netstat -tlnp | grep voidling
```

### Resource Usage

```bash
# CPU and Memory usage
top -p $(pgrep voidling)

# Or using htop
htop -p $(pgrep voidling)
```

## Firewall Configuration

If you have UFW enabled:

```bash
# Allow SSH (if not already allowed)
sudo ufw allow ssh

# Check status
sudo ufw status
```

The bot doesn't need any incoming ports open (it connects outbound to Discord).

## Troubleshooting

### Bot Won't Start

```bash
# Check service status
sudo systemctl status voidling

# Check recent logs
sudo journalctl -u voidling -n 50

# Check permissions
ls -la ~/voidling/voidling
ls -la ~/voidling/data/

# Test run manually
cd ~/voidling
./voidling
```

### Database Issues

```bash
# Check database file exists
ls -la ~/voidling/data/voidling.db

# Check permissions
chmod 644 ~/voidling/data/voidling.db
```

### Out of Memory

If the server runs out of memory, add swap:

```bash
# Create 2GB swap file
sudo fallocate -l 2G /swapfile
sudo chmod 600 /swapfile
sudo mkswap /swapfile
sudo swapon /swapfile

# Make permanent
echo '/swapfile none swap sw 0 0' | sudo tee -a /etc/fstab
```

### Check Discord Token

```bash
# View .env file
cat ~/voidling/.env

# Make sure token is correct and has no extra spaces
```

## Security Best Practices

1. **Use a dedicated user** (not root):
   ```bash
   sudo useradd -m -s /bin/bash voidling
   sudo -u voidling -i
   # Install and configure as voidling user
   ```

2. **Restrict file permissions**:
   ```bash
   chmod 600 ~/.env
   chmod 700 ~/voidling/data
   ```

3. **Keep system updated**:
   ```bash
   sudo apt update
   sudo apt upgrade -y
   ```

4. **Use SSH keys** instead of passwords for server access

5. **Configure fail2ban** to prevent brute force attacks:
   ```bash
   sudo apt install fail2ban -y
   sudo systemctl enable fail2ban
   ```

## Performance Tuning

For better performance on production:

### Increase File Descriptor Limits

Edit `/etc/systemd/system/voidling.service` and add:

```ini
[Service]
LimitNOFILE=65536
```

### Database Optimization

The bot uses SQLite with automatic migrations. For better performance with large databases:

```bash
# Vacuum database periodically (compacts and optimizes)
sqlite3 ~/voidling/data/voidling.db "VACUUM;"

# Add to monthly cron
crontab -e
# Add: 0 0 1 * * sqlite3 /home/your-username/voidling/data/voidling.db "VACUUM;"
```

## Summary of Commands

```bash
# Quick deployment checklist
mkdir -p ~/voidling/data ~/voidling/logs
cp voidling ~/voidling/
chmod +x ~/voidling/voidling
nano ~/voidling/.env  # Add DISCORD_TOKEN
sudo nano /etc/systemd/system/voidling.service  # Create service
sudo systemctl daemon-reload
sudo systemctl enable voidling
sudo systemctl start voidling
sudo systemctl status voidling
```

## Additional Resources

- [systemd Documentation](https://www.freedesktop.org/software/systemd/man/)
- [Go Deployment Best Practices](https://go.dev/doc/code)
- [Discord Bot Best Practices](https://discord.com/developers/docs/topics/community-resources)
