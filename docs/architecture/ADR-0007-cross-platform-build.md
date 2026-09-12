# ADR-0007: Cross-Platform Build & Packaging Pipeline

## Status
Accepted

## Context
While Windows is the primary developer target, AGY++ must support Linux and macOS environments cleanly.

## Decision
1. **GitHub Actions CI/CD**:
   - Matrix builds for:
     - `windows-amd64` (`agy++.exe`)
     - `windows-arm64` (`agy++.exe`)
     - `linux-amd64` (`agy++`)
     - `linux-arm64` (`agy++`)
     - `darwin-amd64` (`agy++`)
     - `darwin-arm64` (`agy++`)
2. **Build Tooling**:
   - Standard Go build flags (`-ldflags="-s -w -X 'github.com/cdrivex4/agy-plus-plus/internal/version.Version=...'"`).
   - Embedding version, git commit, build date, and baseline compatibility version into the binary.
3. **Packaging**:
   - Distribution as standalone zip/tar archives containing the binary, LICENSE, README, and THIRD_PARTY_NOTICES.
   - Dedicated PowerShell installer script (`scripts/install/install.ps1`) for Windows users.

## Consequences
- Single automated workflow produces tested, signed, and hashed release bundles across all platforms.
