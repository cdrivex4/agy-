# ADR-0003: Authentication Delegation & Zero-Secrets Architecture

## Status
Accepted

## Context
Antigravity CLI uses system keyrings and Google OAuth for authentication. AGY++ could theoretically intercept OAuth tokens, scrape keyrings, or proxy requests, but doing so would create massive security vulnerabilities, violate security principles, and risk credential leakage.

## Decision
AGY++ will adopt a **Zero-Secrets Policy**:
1. Never scrape, inspect, or store OAuth tokens, passwords, or API keys.
2. Delegate all authentication lifecycle events (`/login`, `/logout`, `/auth status`) directly to the underlying `agy` command or its documented authentication entrypoints.
3. Automatically sanitize and redact all diagnostics, logs, and trace files to ensure secrets never leak into issue reports or terminal dumps.

## Consequences
- Clean legal and security posture.
- Seamless compatibility with user's existing Google authentication without requiring re-authentication.
