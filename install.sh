#!/bin/bash

# Copyright The OpenTelemetry Authors
# SPDX-License-Identifier: Apache-2.0

set -e

detect() {
    CURRENT_OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    CURRENT_ARCH=$(uname -m)
    if [ "${CURRENT_ARCH}" == "x86_64" ]; then
        CURRENT_ARCH="amd64"
    fi

    echo "Detected platform: ${CURRENT_OS} ${CURRENT_ARCH}"
}

download() {
    DOWNLOAD_URL="https://github.com/alibaba/opentelemetry-go-auto-instrumentation/releases/latest/download/otel-${CURRENT_OS}-${CURRENT_ARCH}"
    EXECUTABLE="otel"

    echo "Downloading from $DOWNLOAD_URL"
    # curl and show progress
    curl -L -o "$EXECUTABLE" "$DOWNLOAD_URL"

    if [ $? -ne 0 ]; then
        echo "Failed to download $DOWNLOAD_URL"
        exit 1
    fi
}

install() {
    INSTALL_DIR="/usr/local/bin"
    if [ ! -f "$EXECUTABLE" ]; then
        echo "Executable $EXECUTABLE not found"
        exit 1
    fi
    echo "Installing $EXECUTABLE to $INSTALL_DIR"
    sudo mv "$EXECUTABLE" "$INSTALL_DIR/"
    sudo chmod +x "$INSTALL_DIR/$EXECUTABLE"

    echo "Installation completed. You can run it using: $INSTALL_DIR/$EXECUTABLE"
}

detect
download
install
