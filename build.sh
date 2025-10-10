#!/bin/bash

APP_NAME="bdterm"
DIR_NAME="release"
# Build for Intel
echo "Building for macOS (Intel)..."
GOOS=darwin GOARCH=amd64 go build -o ${APP_NAME}-amd64

# Build for Apple Silicon
echo "Building for macOS (Apple Silicon)..."
GOOS=darwin GOARCH=arm64 go build -o ${APP_NAME}-arm64

# Create universal binary
echo "Creating universal binary..."
lipo -create -output ${APP_NAME} ${APP_NAME}-amd64 ${APP_NAME}-arm64

# Clean up individual binaries
rm ${APP_NAME}-amd64 ${APP_NAME}-arm64

echo "Universal binary created: ${APP_NAME}"
file ${APP_NAME}

# Optionally, you can strip the binary to reduce its size
strip ${APP_NAME}

rm -rf ${DIR_NAME}
mkdir -p ${DIR_NAME}
mv ${APP_NAME} ${DIR_NAME}/${APP_NAME}
echo "Final binary moved to ${DIR_NAME}/${APP_NAME}"