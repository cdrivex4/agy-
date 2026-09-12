# ADR-0002: Upstream AGY Integration and Adapter Boundary

## Status
Accepted

## Context
Upstream Antigravity CLI evolves rapidly with frequent releases. Directly coupling AGY++ internals to specific AGY versions or hard-coded flags would lead to fragility and frequent breaking bugs on AGY updates.

## Decision
All interactions with upstream AGY must be isolated behind strict adapter interfaces located in `internal/upstream/agy/`:
- `AgyProcessAdapter`: Manages execution of the `agy` binary in interactive and stream-json modes.
- `AgyCapabilityDetector`: Probes installed AGY binary capabilities at runtime.
- `AgyConfigAdapter`: Reads upstream projects, workspaces, and configurations without modifying them.
- `AgyAuthAdapter`: Detects authentication status and triggers official login flows.

Direct version checks (`if version >= ...`) are prohibited outside the adapter package; core application logic queries capability flags (e.g. `capabilities.Supports("stream-json")`).

## Consequences
- Upstream changes only require updates inside `internal/upstream/agy/`.
- Regression tests and mock adapters can be built easily using sanitized test fixtures.
