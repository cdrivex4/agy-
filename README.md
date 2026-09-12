# AGY++ (agy++)

> **Antigravity CLI on steroids** — an independent, extensible companion application for Google Antigravity CLI, combining the power of the AGY agent engine with the power-user experience of Gemini CLI.

---

## What is AGY++?
AGY++ wraps your existing Antigravity CLI (`agy`) installation to provide:
- ⚡ **True YOLO Mode**: Run with auto-approved permissions (`--yolo`, `-y`, `/yolo`) while preserving strict security guardrails against prompt injection and maintaining sandbox boundaries.
- 📋 **Plan Mode & Guardrails**: Enforce read-only planning workflows before making changes.
- 🔄 **Session Checkpointing & Compacting**: Bookmark, rewind, and compact conversations without destroying raw transcripts.
- 🛡️ **Zero-Secrets Security**: Transparent delegation to your system keyring and Google Sign-In; no credential scraping or leaking.
- 🧩 **100% Upstream Compatible**: Seamlessly coexists with your current AGY installation, MCP servers, custom skills, and plugins.

---

## Quick Start

### Prerequisites
- Windows 10/11, macOS, or Linux
- Google Antigravity CLI (`agy`) installed and authenticated

### Installation
Download the latest binary release from GitHub Releases or build from source:

```powershell
# Build from source
go build -o agy++.exe ./cmd/agy-plus-plus
```

### Basic Usage
```bash
# Launch interactive AGY++ shell
agy++

# Run with YOLO auto-approval mode
agy++ --yolo

# Run non-interactively with a prompt
agy++ -p "Review this codebase and summarize architecture"

# Check upstream compatibility and system status
agy++ diagnostics
```

---

## Architecture & Compatibility
AGY++ strictly isolates all interactions with the upstream `agy` binary behind clean-room adapters in `internal/upstream/agy/`.
- Upstream Target: **Antigravity CLI v1.2.2**
- See [ARCHITECTURE.md](ARCHITECTURE.md) and [UPSTREAM_COMPATIBILITY.md](UPSTREAM_COMPATIBILITY.md) for architectural details.

---

## License
Apache License 2.0. See [LICENSE](LICENSE) for details.
