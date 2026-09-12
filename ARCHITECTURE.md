# AGY++ Architecture

## 1. High-Level Design

AGY++ is designed as an independent CLI/TUI companion application that enriches Antigravity CLI while maintaining full compatibility with the upstream ecosystem.

```
                    +------------------------------------+
                    |               AGY++                |
                    +-----------------+------------------+
                                      |
              +-----------------------+-----------------------+
              |                                               |
              v                                               v
     [AGY++ UX & Features]                        [AGY Compatibility Layer]
  - YOLO Approval Mode                          - Discovery & Version Check
  - Plan Mode Guardrails                        - Process Execution (TUI/JSON)
  - Persistent Slash Commands                   - Session & Project Reader
  - Session Checkpoint & Compact                - Auth & Keyring Delegation
  - Diagnostic & Upstream Auditor               - MCP & Plugin Integration
              |                                               |
              +-----------------------+-----------------------+
                                      |
                                      v
                        +---------------------------+
                        |  Google Antigravity (AGY) |
                        |     Native Binary         |
                        +---------------------------+
```

## 2. Core Principles
1. **Upstream Isolation**: All upstream AGY communication is contained within `internal/upstream/agy/`. No business logic directly queries upstream version strings or manipulates upstream data files.
2. **Read-Only Upstream Storage**: Upstream SQLite conversation databases and configuration caches are treated as strictly read-only.
3. **Companion Data Store**: AGY++ maintains its own SQLite database in `%APPDATA%\AgyPlusPlus\sessions.db` (or `~/.config/agy-plus-plus/sessions.db`) for bookmarks, checkpoints, and extended context metadata.
4. **Security by Design**: Zero credential scraping or logging. YOLO activation is restricted strictly to human input.

## 3. Directory Layout
- `cmd/agy-plus-plus/`: Entry point for `agy++.exe`.
- `internal/app/`: Application lifecycle and initialization.
- `internal/cli/`: CLI flag definitions, subcommands, and dispatching.
- `internal/tui/`: Interactive terminal UI and renderers.
- `internal/commands/`: Slash command registry (`/yolo`, `/status`, `/compact`, etc.).
- `internal/policy/`: YOLO and security permission gatekeepers.
- `internal/upstream/agy/`: Adapters, process supervision, and capability detectors for AGY.
- `internal/session/`: Dual-layer session management.
- `internal/diagnostics/`: System health, redaction engine, and diagnostic dumps.
- `tests/fixtures/agy/`: Sanitized upstream output fixtures for automated testing.
