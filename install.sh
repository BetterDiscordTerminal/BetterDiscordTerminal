#!/bin/bash

# BetterDiscord Terminal Installer
# Author: nmn
# Website: https://pand.dev
# Source: https://github.com/BetterDiscordTerminal/BetterDiscordTerminal
# License: Apache-2.0

# This script downloads and verifies the bdterm installer binary and runs it.
# It is recommended to review the script before executing it.
# Usage: curl -fsSL https://raw.githubusercontent.com/BetterDiscordTerminal/BetterDiscordTerminal/main/install.sh | bash

set -e

# Colors
RED='\033[0;31m'
YELLOW='\033[0;33m'
GREEN='\033[0;32m'
NC='\033[0m'

# Configuration
BINARY_URL="https://github.com/BetterDiscordTerminal/BetterDiscordTerminal/releases/download/v0.0.1/bdterm"
BINARY_NAME="bdterm"
TMP_DIR="/tmp"
BINARY_PATH="${TMP_DIR}/${BINARY_NAME}"
EXPECTED_SHA256="ae7d7d450aebd508dd66be7f611ac83c9c93debbc7ee247af9f5d57f01732a5e"

# Detect architecture
ARCH=$(uname -m)
if [[ "$ARCH" != "arm64" ]]; then
    echo -e "${RED}Error: Unsupported architecture: $ARCH${NC}"
    echo -e "${YELLOW}BDTerm only supports Apple Silicon (ARM64) Macs.${NC}"
    echo -e "${YELLOW}Intel Macs are not supported in this version.${NC}"
    exit 1
fi

# Download binary
if ! curl -fsSL "${BINARY_URL}" -o "${BINARY_PATH}"; then
    echo -e "${RED}Error: Failed to download installer${NC}"
    exit 1
fi

# Verify SHA-256 checksum
ACTUAL_SHA256=$(shasum -a 256 "${BINARY_PATH}" | awk '{print $1}')
if [[ "$ACTUAL_SHA256" != "$EXPECTED_SHA256" ]]; then
    echo -e "${RED}Error: Checksum verification failed${NC}"
    echo -e "${RED}  Expected: $EXPECTED_SHA256${NC}"
    echo -e "${RED}  Got: $ACTUAL_SHA256${NC}"
    echo -e "${YELLOW}Installer may be corrupted/tampered with or the expected checksum is outdated.${NC}"
    echo -e "${YELLOW}The program was not run - your system is safe.${NC}"
    rm -f "${BINARY_PATH}"
    exit 1
fi

echo -e "${GREEN}Checksum verification passed${NC}"

# Make executable
chmod +x "${BINARY_PATH}"

# Run the binary
"${BINARY_PATH}"

# Clean up
rm -f "${BINARY_PATH}"
