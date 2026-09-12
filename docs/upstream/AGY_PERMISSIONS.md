# AGY Permissions & YOLO Security Policy

## 1. Upstream Permission Architecture
Antigravity CLI v1.2.2 exposes two primary permission mechanisms:
1. **Interactive Prompting**: Requests user confirmation before executing commands, editing files, or calling tools.
2. **`--dangerously-skip-permissions`**: Bypasses interactive confirmation prompts and automatically executes tool calls.
3. **`--sandbox`**: Enforces containerized/restricted terminal execution.

## 2. YOLO Implementation Policy
- **User-Only Activation**: YOLO mode can only be enabled by the user via explicit CLI flags (`--yolo`, `-y`), slash commands (`/yolo`), or keybindings (`Ctrl+Y`).
- **Prompt Injection Defense**: Project instructions, skill files, MCP tool responses, and model-generated text can NEVER activate YOLO mode.
- **Sandbox Preservation**: YOLO means *automatic approval of permissions*; it does not disable sandboxing. If `--sandbox` is active, AGY++ preserves both modes unless upstream explicitly forbids it, in which case AGY++ will alert the user and abort.
