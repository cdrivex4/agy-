# AGY Session Architecture & AGY++ Metadata Layer

## 1. Upstream Session Storage
- Antigravity CLI v1.0.4+ transitioned conversation storage to SQLite `.db` databases.
- Conversations can be resumed via `--continue` (`-c`) or `--conversation <id>`.
- Project-to-workspace mappings are indexed in `~/.gemini/antigravity-cli/cache/projects.json`.

## 2. AGY++ Stance: Read-Only Upstream Access
- Upstream SQLite databases are treated strictly as **read-only compatibility data**.
- AGY++ will never execute `INSERT`, `UPDATE`, or `DELETE` statements on AGY's internal databases.

## 3. AGY++ Companion Session Store
AGY++ maintains its own SQLite database in `%APPDATA%\AgyPlusPlus\sessions.db`:
- Stores bookmarks, user notes, checkpoint references, compacted context snapshots, and tags.
- Links records to upstream sessions via `agy_session_id`.
- Provides transactional migrations, automatic backups, and rollback safeguards.
