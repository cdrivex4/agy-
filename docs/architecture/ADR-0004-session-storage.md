# ADR-0004: Dual Session Architecture

## Status
Accepted

## Context
Upstream AGY stores conversation trajectories in internal SQLite databases. Writing directly to AGY's SQLite database risks corrupting upstream state or conflicting when upstream updates its internal database schema.

## Decision
AGY++ will employ a **Dual Session Architecture**:
1. **Upstream SQLite Database**: Treated as strictly read-only. Used solely for discovering past conversation IDs and reading history when resuming.
2. **AGY++ Companion Database**: Stored in `%APPDATA%\AgyPlusPlus\sessions.db` (or `~/.config/agy-plus-plus/sessions.db` on Unix). Contains AGY++ specific metadata: bookmarks, tags, checkpoint markers, compact representations, and user notes, keyed by `agy_session_id`.

## Consequences
- Upstream stability is completely protected from database lock conflicts and schema drift.
- Full independence for rich AGY++ features (checkpoints, bookmarks, compacting).
