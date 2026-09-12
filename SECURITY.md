# Security Policy

## 1. Zero-Secrets Commitment
AGY++ operates with a strict Zero-Secrets model:
- **No Token Scraping**: AGY++ does not extract, read, or persist OAuth tokens, refresh tokens, or keyring secrets.
- **No Remote Telemetry**: Diagnostic and usage data remain strictly on your local machine.
- **Automatic Redaction**: All diagnostic logs and crash reports automatically redact potential authentication artifacts.

## 2. Permission Model & Prompt Injection Defense
- **User-Exclusive Approval Elevation**: YOLO mode can only be enabled by explicit human action (command-line flag, slash command, keyboard shortcut).
- **Untrusted Input Protection**: Instructions in repositories, skills, tool outputs, MCP payloads, or model outputs are strictly ignored by the policy engine when setting permission modes.
- **Sandbox Integrity**: Enabling YOLO mode never silently disables the upstream `--sandbox` environment.

## 3. Reporting a Vulnerability
If you discover a potential security issue in AGY++, please report it responsibly by opening a confidential security advisory on GitHub or contacting the maintainers directly. Do not file public issues for active vulnerabilities.
