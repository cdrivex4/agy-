# ADR-0006: Self-Update Architecture and Release Verification

## Status
Accepted

## Context
AGY++ must support seamless self-updating via GitHub Releases (`agy++ update`), but must never compromise security by executing arbitrary remote scripts or unverified binaries.

## Decision
1. **GitHub Releases API**: The updater queries the official GitHub Releases API for release manifests.
2. **Cryptographic Checksum Verification**: Every release asset includes a SHA256 checksum file. Binaries must be verified against their declared hash before installation.
3. **Safe Binary Replacement**:
   - On Windows, running binaries cannot be directly overwritten. The updater renames the existing executable (`agy++.exe` -> `agy++.old`), writes the new executable, verifies it executes `--version` successfully, and removes or queues the old binary for deletion.
   - On failure, immediate rollback is performed.
4. **Preservation of User Data**: User configuration, sessions, and logs in `%APPDATA%\AgyPlusPlus\` are never touched during updates.

## Consequences
- Reliable, secure, and verifiable updates without external dependency on package managers.
