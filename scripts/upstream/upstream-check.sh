#!/usr/bin/env bash
set -euo pipefail

echo "=== AGY++ Upstream Version Auditor ==="

if ! command -v agy &> /dev/null; then
    echo "[!] agy command not found in PATH."
    exit 1
fi

INSTALLED_VERSION=$(agy --version 2>&1 | tr -d '\r')
echo "Found installed AGY: ${INSTALLED_VERSION}"

echo ""
echo "Compatibility Status:"
echo "  Installed Version:  ${INSTALLED_VERSION}"
echo "  Tested Baseline:    1.2.2"

if [ "${INSTALLED_VERSION}" = "1.2.2" ]; then
    echo "  Status:             VERIFIED COMPATIBLE"
else
    echo "  Status:             VERSION MISMATCH - Review changelog and update fixtures."
fi
