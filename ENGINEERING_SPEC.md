# AGY++ Engineering Specification

## 1. Scope & Objective
This specification defines the functional, technical, and operational boundaries of **AGY++**. AGY++ restores power-user features inspired by Gemini CLI (YOLO mode, enhanced interactive commands, session checkpoints) over the production Antigravity CLI engine.

## 2. Technical Stack
- **Language**: Go 1.23+
- **Primary Binary**: `agy++.exe` (Windows x86_64, ARM64)
- **Secondary Targets**: Linux, macOS
- **State Storage**: SQLite for AGY++ metadata
- **Upstream Engine Target**: Google Antigravity CLI v1.2.2+

## 3. Subsystem Specifications

### 3.1 Upstream Adapter Subsystem (`internal/upstream/agy`)
- **Discovery**: Automatically locate `agy` binary via system `PATH`, standard installation directories (`%LOCALAPPDATA%\agy\bin\agy.exe`), or explicit environment overrides (`AGY_BINARY_PATH`).
- **Capability Detection**: Probe flags and subcommands without user impact.
- **Process Supervision**: Launch upstream AGY in interactive TUI or headless (`--print`, `stream-json`) mode with graceful signal handling and cleanup.

### 3.2 Policy Subsystem (`internal/policy`)
- Enforces strict human-only activation for YOLO mode (`-y`, `--yolo`, `/yolo`, `Ctrl+Y`).
- Rejects permission elevation from project instructions, tools, MCP responses, or prompt text.
- Honors `--sandbox` isolation boundary.

### 3.3 Command Registry (`internal/commands`)
- Implements slash command dispatch (`/help`, `/status`, `/login`, `/logout`, `/yolo`, `/plan`, `/compact`, `/checkpoint`, `/diff`, `/quit`).
- Maps commands to required backend capabilities and provides graceful degradation with clear user feedback when upstream lacks support.

### 3.4 Diagnostics & Redaction Subsystem (`internal/diagnostics`)
- Generates system and session diagnostics with automatic regex-based redaction of auth tokens, API keys, cookies, and sensitive project paths.
