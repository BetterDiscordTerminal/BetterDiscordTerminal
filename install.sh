#!/bin/bash

# BetterDiscord Terminal Installer
# Author: nmn
# Website: https://pand.dev
# Source: https://github.com/pandeynmn/bdterm
# License: Apache-2.0

# This script downloads and verifies the bdterm installer binary and runs it.
# It is recommended to review the script before executing it.
# Usage: curl -fsSL https://raw.githubusercontent.com/pandeynmn/bdterm/main/install.sh | bash
#        curl -fsSL https://bd.pand.dev/install.sh | bash

set -e

# Colors
RED='\033[0;31m'
YELLOW='\033[0;33m'
GREEN='\033[0;32m'
NC='\033[0m'

# Configuration
BINARY_URL="https://github.com/BetterDiscordTerminal/BetterDiscordTerminal/releases/download/v1.0.2/bdterm"
BINARY_NAME="bdterm"
TMP_DIR="/tmp"
BINARY_PATH="${TMP_DIR}/${BINARY_NAME}"
EXPECTED_SHA256="2ac4eb8802b50f8e3e7610ec09d4836af46becceca63de24eef3dda7d5113dac"

# Detect architecture
ARCH=$(uname -m)
if [[ "$ARCH" != "x86_64" && "$ARCH" != "arm64" ]]; then
    echo -e "${RED}Error: Unsupported architecture: $ARCH${NC}"
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
