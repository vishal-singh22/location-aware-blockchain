#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e

echo "Starting the installation process..."

# Step 1: Update the system packages
echo "Updating system packages..."
sudo apt update && sudo apt upgrade -y

# Step 2: Install required dependencies
echo "Installing required dependencies..."
sudo apt install -y curl git build-essential docker.io docker-compose

# Step 3: Install Go (Golang)
echo "Installing Go..."
GO_VERSION="1.20.6"  # Specify the required Go version
wget https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz
rm go${GO_VERSION}.linux-amd64.tar.gz

# Add Go to PATH (if not already in .bashrc or .zshrc)
if ! grep -q "/usr/local/go/bin" ~/.bashrc; then
    echo "export PATH=\$PATH:/usr/local/go/bin" >> ~/.bashrc
fi
source ~/.bashrc

# Verify Go installation
echo "Verifying Go installation..."
go version

# Step 4: Install Hyperledger Fabric prerequisites
echo "Installing Hyperledger Fabric prerequisites..."
sudo apt install -y jq

# Step 5: Pull Docker images for Hyperledger Fabric
echo "Pulling Hyperledger Fabric Docker images..."
curl -sSL https://bit.ly/2ysbOFE | bash -s

# Step 6: Verify Docker and Docker Compose installation
echo "Verifying Docker and Docker Compose installation..."
docker --version
docker-compose --version

# Final message
echo "Installation complete. Please restart your terminal session to apply all changes."
