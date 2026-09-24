#!/bin/bash

# VecScan Installation Script
# Run: chmod +x install.sh && sudo ./install.sh

set -e

echo "╔═══════════════════════════════════════════════════════════╗"
echo "║  VecScan by Vectalith Labs - Installation Script          ║"
echo "╚═══════════════════════════════════════════════════════════╝"
echo

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Check if running as root for some operations
if [ "$EUID" -ne 0 ]; then
  echo -e "${YELLOW}Note: Some features require root privileges${NC}"
fi

# Check for Go
if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is not installed${NC}"
    echo "Please install Go 1.21 or later: https://golang.org/doc/install"
    exit 1
fi

GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
echo -e "${GREEN}✓ Go version: $GO_VERSION${NC}"

# Install dependencies
echo "[*] Installing dependencies..."
go mod download
go mod tidy

# Build
echo "[*] Building VecScan..."
make build

# Install binary
echo "[*] Installing binary..."
sudo cp build/vecscan /usr/local/bin/
sudo chmod +x /usr/local/bin/vecscan

# Verify installation
if command -v vecscan &> /dev/null; then
    echo -e "${GREEN}✓ VecScan installed successfully!${NC}"
    vecscan --version 2>/dev/null || true
    echo
    echo "Usage examples:"
    echo "  vecscan 192.168.1.1                    # Basic scan"
    echo "  vecscan 192.168.1.1 -p 1-65535         # Full port scan"
    echo "  vecscan 192.168.1.1 --syn -rate 500    # SYN scan, 500 pps"
    echo "  vecscan 192.168.1.1 -o results.json    # JSON output"
else
    echo -e "${RED}Installation failed${NC}"
    exit 1
fi