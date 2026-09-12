# ADR-0005: YOLO Mode Security Model and Injection Guardrails

## Status
Accepted

## Context
YOLO mode allows automatic approval of tool actions. If an untrusted input (such as repository instructions, an MCP server output, or prompt injection) could trigger YOLO mode, an attacker could achieve arbitrary code execution on the user's workstation.

## Decision
1. **User-Exclusive Activation**: YOLO mode can only be engaged by explicit human action:
   - Command-line argument: `--yolo`, `-y`, `--approval-mode=yolo`
   - Interactive TUI command: `/yolo`
   - Optional interactive shortcut: `Ctrl+Y`
2. **Untrusted Input Isolation**: Neither project instructions, skill files, system prompts, MCP responses, nor LLM model outputs can switch the approval mode to YOLO.
3. **Persistent Visual Indication**: Whenever YOLO is active, a prominent indicator (`MODE: YOLO`) must be displayed in the status bar and terminal banner.
4. **Sandbox Preservation**: YOLO mode auto-approves permissions; it does not disable sandboxing. If `--sandbox` is enabled, AGY++ preserves both modes.

## Consequences
- Protects users against prompt injection attacks attempting to grant themselves unprompted tool permissions.
