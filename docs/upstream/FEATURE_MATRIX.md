# AGY vs. Gemini CLI vs. AGY++ Feature Matrix

This document benchmarks capabilities across Antigravity CLI (upstream), Gemini CLI (reference), and AGY++ (target product).

| Capability | Antigravity CLI (1.2.2) | Gemini CLI (Reference) | AGY++ (Target) | AGY++ Implementation Strategy |
| :--- | :--- | :--- | :--- | :--- |
| **YOLO Mode** | `--dangerously-skip-permissions` flag | `-y`, `/yolo`, auto-approve | `--yolo`, `-y`, `/yolo`, `Ctrl+Y` | Wrap AGY flag; ensure user-only activation; prominent indicator. |
| **Plan Mode** | `--mode plan` flag | `/plan`, read-only inspection | `--plan`, `/plan` | Integrate with AGY `--mode plan` + enforced local read-only guardrails. |
| **Compact** | Internal context management | `/compact` command | `/compact` | Preserve raw transcript; create compacted summary layer in AGY++ store. |
| **Rewind / Checkpoints** | Not exposed in CLI | `/rewind`, checkpoints | `/checkpoint`, `/rewind`, `/diff` | Git commit/tree checkpoint + conversation checkpoint in `~/.agy-plus-plus/`. |
| **Auth Delegation** | System keyring / Google OAuth | System keyring / OAuth | `/login`, `/logout`, `/auth status` | Transparent delegation to AGY without credential interception. |
| **Headless Output** | `--print` (`text`, `json`, `stream-json`) | Headless print | `agy++ -p` with streaming NDJSON | Wrap AGY headless runner; structured error reporting and cancellation. |
| **Session Metadata** | SQLite `.db` conversation store | SQLite store | Augmenting SQLite store | Read-only access to AGY sessions; companion AGY++ session database. |
| **Upstream Tracking** | N/A | N/A | Automated update checks & diffs | `scripts/upstream-check.ps1`, `scripts/upstream-diff.ps1`, change classification. |
| **Windows Support** | Single native `agy.exe` | Multi-platform CLI | First-class Windows `agy++.exe` | Go cross-compilation; PowerShell installer; Windows AppData conventions. |
