#!/usr/bin/env bash
set -euo pipefail

echo "Fetching upstream changelog..."
REPORT_DATE=$(date +%Y-%m-%d)
REPORT_DIR="$(dirname "$0")/../../docs/upstream/reports"
mkdir -p "${REPORT_DIR}"
REPORT_FILE="${REPORT_DIR}/${REPORT_DATE}-agy-update.md"

cat <<EOF > "${REPORT_FILE}"
# Upstream Audit Report - ${REPORT_DATE}

## Installed AGY Version
$(agy --version 2>&1)

## Upstream Changelog Summary
$(agy changelog 2>&1)
EOF

echo "Saved upstream report to: ${REPORT_FILE}"
