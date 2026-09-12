#!/usr/bin/env bash
set -euo pipefail

OUTPUT="${1:-agy++}"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
VERSION=$(cat "${SCRIPT_DIR}/../../VERSION" | tr -d '\r\n')
COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE=$(date +%Y-%m-%d)

LDFLAGS="-s -w -X 'github.com/cdrivex4/agy-plus-plus/internal/version.Version=${VERSION}' \
-X 'github.com/cdrivex4/agy-plus-plus/internal/version.GitCommit=${COMMIT}' \
-X 'github.com/cdrivex4/agy-plus-plus/internal/version.BuildDate=${DATE}' \
-X 'github.com/cdrivex4/agy-plus-plus/internal/version.BaselineAgy=1.2.2'"

echo "Building AGY++ v${VERSION} (${COMMIT})..."
go build -ldflags "${LDFLAGS}" -o "${OUTPUT}" ./cmd/agy-plus-plus
echo "Build complete: ${OUTPUT}"
