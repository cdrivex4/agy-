# AGY Configuration & Settings Precedence

## 1. Upstream Configuration Layout
Antigravity CLI uses the following storage locations:
- Cache & Projects: `~/.gemini/antigravity-cli/cache/projects.json`
- Plugins & Custom Skills: `~/.gemini/config/`
- MCP Servers: `~/.gemini/config/mcp_config.json`
- Environment Variables:
  - `AGY_CLI_DISABLE_LATEX`: Toggles terminal LaTeX rendering.
  - `AGY_CLI_HIDE_ACCOUNT_INFO`: Suppresses user email and account tier in terminal headers.

## 2. AGY++ Effective Configuration Hierarchy
AGY++ implements a layered configuration system (`EffectiveConfig`):
1. Command-line flags (highest precedence)
2. Project-level AGY++ configuration (`.agypp/config.yaml`)
3. User-level AGY++ configuration (`%APPDATA%\AgyPlusPlus\config.yaml`)
4. Upstream AGY configuration (`~/.gemini/...`)
5. Environment variables and built-in defaults (lowest precedence)

AGY++ never writes to or alters upstream configuration files without explicit user command, and preserves full reversibility.
