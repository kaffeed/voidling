#!/bin/bash
# Server Setup Script for Voidling Discord Bot
# Run this script on your Ubuntu server to prepare for deployment

set -e

# Configuration (modify these as needed)
APP_USER="voidling"
DEPLOY_PATH="/opt/voidling"
SERVICE_NAME="voidling"

echo "=========================================="
echo "Voidling Discord Bot - Server Setup"
echo "=========================================="
echo ""

# Check if running as root or with sudo
if [ "$EUID" -ne 0 ]; then
    echo "Please run this script with sudo:"
    echo "sudo bash setup-server.sh"
    exit 1
fi

# Step 1: Create application user
echo "[1/6] Creating application user..."
if id "$APP_USER" &>/dev/null; then
    echo "User $APP_USER already exists"
else
    adduser --system --group --home "$DEPLOY_PATH" --shell /bin/bash "$APP_USER"
    echo "User $APP_USER created"
fi

# Step 2: Create deployment directory
echo "[2/6] Creating deployment directory..."
mkdir -p "$DEPLOY_PATH/data"
chown -R "$APP_USER:$APP_USER" "$DEPLOY_PATH"
chmod 755 "$DEPLOY_PATH"
echo "Deployment directory created at $DEPLOY_PATH"

# Step 3: Install required packages
echo "[3/6] Installing required packages..."
apt-get update
apt-get install -y sqlite3
echo "Required packages installed"

# Step 4: Create .env file template
echo "[4/6] Creating .env file template..."
if [ ! -f "$DEPLOY_PATH/.env" ]; then
    cat > "$DEPLOY_PATH/.env" << 'EOF'
# Discord Bot Configuration
DISCORD_TOKEN=your_bot_token_here

# Database Configuration
DATABASE_PATH=/opt/voidling/data/voidling.db

# Logging
LOG_LEVEL=info

# Discord Guild ID (optional, for faster development command registration)
# DISCORD_GUILD_ID=your_guild_id_here

# Coordinator Role ID (optional)
# COORDINATOR_ROLE_ID=your_role_id_here
EOF
    chown "$APP_USER:$APP_USER" "$DEPLOY_PATH/.env"
    chmod 600 "$DEPLOY_PATH/.env"
    echo ".env file created - PLEASE EDIT IT WITH YOUR ACTUAL VALUES"
    echo "Location: $DEPLOY_PATH/.env"
else
    echo ".env file already exists, skipping"
fi

# Step 5: Setup systemd service
echo "[5/6] Setting up systemd service..."
if [ -f "voidling.service" ]; then
    cp voidling.service /etc/systemd/system/
elif [ -f "../scripts/voidling.service" ]; then
    cp ../scripts/voidling.service /etc/systemd/system/
else
    echo "Warning: voidling.service file not found"
    echo "You'll need to manually copy it to /etc/systemd/system/"
fi

# Configure sudo permissions for service management (optional but recommended)
if [ ! -f "/etc/sudoers.d/$APP_USER" ]; then
    cat > "/etc/sudoers.d/$APP_USER" << EOF
# Allow $APP_USER to manage voidling service without password
$APP_USER ALL=(ALL) NOPASSWD: /bin/systemctl start $SERVICE_NAME
$APP_USER ALL=(ALL) NOPASSWD: /bin/systemctl stop $SERVICE_NAME
$APP_USER ALL=(ALL) NOPASSWD: /bin/systemctl restart $SERVICE_NAME
$APP_USER ALL=(ALL) NOPASSWD: /bin/systemctl status $SERVICE_NAME
$APP_USER ALL=(ALL) NOPASSWD: /bin/journalctl -u $SERVICE_NAME*
EOF
    chmod 440 "/etc/sudoers.d/$APP_USER"
    echo "Sudo permissions configured for service management"
fi

systemctl daemon-reload
systemctl enable "$SERVICE_NAME"
echo "Systemd service configured"

# Step 6: Setup SSH authorized_keys (manual step reminder)
echo "[6/6] SSH Key Setup..."
AUTHORIZED_KEYS="$DEPLOY_PATH/.ssh/authorized_keys"
mkdir -p "$DEPLOY_PATH/.ssh"
touch "$AUTHORIZED_KEYS"
chown -R "$APP_USER:$APP_USER" "$DEPLOY_PATH/.ssh"
chmod 700 "$DEPLOY_PATH/.ssh"
chmod 600 "$AUTHORIZED_KEYS"
echo "SSH directory created at $DEPLOY_PATH/.ssh"

echo ""
echo "=========================================="
echo "Setup Complete!"
echo "=========================================="
echo ""
echo "Next Steps:"
echo ""
echo "1. Edit the .env file with your Discord bot token:"
echo "   sudo nano $DEPLOY_PATH/.env"
echo ""
echo "2. Add your GitHub Actions SSH public key to:"
echo "   $AUTHORIZED_KEYS"
echo "   Example:"
echo "   echo 'ssh-ed25519 AAAAC3Nza...' | sudo tee -a $AUTHORIZED_KEYS"
echo ""
echo "3. Test SSH connection:"
echo "   ssh -i /path/to/private_key $APP_USER@$(hostname -I | awk '{print $1}')"
echo ""
echo "4. Configure GitHub secrets (see implement/SECRETS.md):"
echo "   - SSH_HOST: $(hostname -I | awk '{print $1}')"
echo "   - SSH_USERNAME: $APP_USER"
echo "   - SSH_KEY: <your private key>"
echo "   - SSH_PORT: 22"
echo "   - DEPLOY_PATH: $DEPLOY_PATH"
echo ""
echo "5. Push a version tag to trigger deployment:"
echo "   git tag v1.0.0 && git push origin v1.0.0"
echo ""
echo "Useful Commands:"
echo "  sudo systemctl start $SERVICE_NAME    # Start the bot"
echo "  sudo systemctl stop $SERVICE_NAME     # Stop the bot"
echo "  sudo systemctl status $SERVICE_NAME   # Check status"
echo "  sudo journalctl -u $SERVICE_NAME -f   # View logs"
echo ""
