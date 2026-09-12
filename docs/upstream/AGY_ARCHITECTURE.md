# Upstream AGY Architecture Analysis

## 1. Overview
Antigravity CLI (AGY) is Google's terminal interface to the core agent engine used by Antigravity 2.0. In version 1.2.2, it is distributed as a single native executable (e.g. `agy.exe` on Windows).

## 2. Key Subsystems Identified

### Execution Modes
- **Interactive TUI**: Launched by default or with `--prompt-interactive` (`-i`).
- **Non-Interactive Print**: Launched with `--print` (`-p`) or `--prompt`. Supports `--output-format` (`text`, `json`, `stream-json`) and `--input-format` (`text`, `stream-json`).
- **Agent Modes**:
  - `accept-edits`: Default operational mode.
  - `plan`: Read-only/planning mode, preventing unintended modifications.
- **Permission Modes**:
  - Standard interactive confirmation prompts.
  - `--dangerously-skip-permissions`: Auto-approves tool actions without interactive prompting (upstream foundation for YOLO mode).
  - Sandbox mode: `--sandbox` enables terminal execution restrictions.

### Subcommands
- `models`: Lists available foundation models (Gemini, Claude, GPT-OSS).
- `agents`: Lists specialized agents.
- `mcp`: Server lifecycle management (`add`, `remove`, `list`, `enable`, `disable`).
- `plugin` / `plugins`: Plugin discovery and lifecycle management.
- `changelog`: Releases and version updates.
- `update`: Self-update pipeline.

### State & Storage
- Central cache: `~/.gemini/antigravity-cli/cache/projects.json` tracks project workspaces.
- Conversation storage: SQLite `.db` databases (shared with Antigravity 2.0).
- Configuration: `~/.gemini/config/` for plugins and MCP configurations.
