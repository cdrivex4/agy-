# Upstream Repositories and Compatibility Tracking

This document tracks upstream repositories, observed versions, and compatibility status for **AGY++**.

## 1. Upstream Targets

### Google Antigravity CLI
- **Repository**: https://github.com/google-antigravity/antigravity-cli
- **Branch**: main
- **Observed Version**: 1.2.2 (detected locally on Windows)
- **Observed Date**: 2026-09-12
- **Documentation**: https://antigravity.google/docs/cli/overview

### Google Gemini CLI (Reference Architecture)
- **Repository**: https://github.com/google-gemini/gemini-cli
- **Branch**: main
- **Observed Date**: 2026-09-12
- **Reference**: Historical CLI UX, session checkpointing, approval models, compact/rewind.

---

## 2. Compatibility Matrix

| AGY Version | Tested | Compatible | Notes | Required Adapter Changes |
| :--- | :--- | :--- | :--- | :--- |
| **1.2.2** | Yes | Yes | Baseline installed version on Windows. Supports --dangerously-skip-permissions, --mode plan, --prompt-interactive, stream-json, SQLite conversation storage. | Baseline adapter targets v1.2.2 flags and commands. |
| **1.1.16** | No | Expected | Historical reference baseline from original project charter. | Legacy flag mapping if invoked on older installs. |

---

## 3. Upstream Compatibility Principles

1. **Clean-Room Adapter**: All AGY-specific flags, outputs, and protocols reside strictly behind adapters in internal/upstream/agy/.
2. **Never Vendor Upstream**: Upstream source code is never wholesale vendored or checked in.
3. **Never Modify Upstream**: User's local AGY binary and installation files are never patched or modified.
4. **No Assumption Sprawl**: Code outside adapters checks capabilities.Supports(...) rather than checking specific version strings.
