# Contributing to AGY++

Thank you for your interest in contributing to AGY++!

## Contribution Guidelines
1. **Architecture Boundaries**: Keep upstream AGY integration code confined within `internal/upstream/agy/`.
2. **Zero-Secrets Policy**: Never add code that scrapes, intercepts, or logs authentication tokens.
3. **Prompt Injection Guardrails**: Ensure that untrusted model responses or repo contents can never alter security policy or activate YOLO mode.
4. **Testing**: Add unit tests or compatibility fixtures for every new feature or upstream adapter change.
5. **Git Commits**: Use conventional commits (`feat:`, `fix:`, `docs:`, `test:`, `chore:`).
