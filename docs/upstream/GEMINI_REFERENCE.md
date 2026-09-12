# Gemini CLI Reference Architecture & Provenance

## 1. Role of Gemini CLI in AGY++
Gemini CLI serves as the historical UX and workflow reference for AGY++. Key capabilities from Gemini CLI that AGY++ aims to revive and modernize:
- Streamlined interactive shell with persistent prompt.
- Intuitive slash command framework (`/compact`, `/rewind`, `/diff`, `/yolo`).
- Clean visual status indicators (active model, approval mode, project workspace).
- Graceful session interruption and state resumption.

## 2. Provenance and Code Isolation
- No blind copying of Gemini CLI code.
- Where reference designs or algorithms (e.g. diff viewing, terminal styling) are adapted from open-source Gemini CLI repositories, their licenses must be verified, proper notices preserved, and logged in `THIRD_PARTY_NOTICES.md`.
