# Upstream Compatibility Guide

## Policy Overview
AGY++ wraps and complements Google Antigravity CLI without forking or patching upstream code. Maintaining continuous compatibility as Antigravity releases new versions is a tier-1 project requirement.

## Compatibility Verification Process
When Google releases a new version of Antigravity CLI:
1. **Detect**: Run `./scripts/upstream/upstream-check.ps1` to detect the new upstream version and changelog.
2. **Inspect & Classify**: Review the upstream diff and classify the changes:
   - `NONE`: No impact on AGY++ adapters.
   - `INFORMATIONAL`: Upstream documentation or cosmetic changes.
   - `COMPATIBILITY_RELEVANT`: Changes to flags, session schemas, or output formats.
   - `FEATURE_OPPORTUNITY`: New upstream capabilities to expose.
   - `BREAKING`: Removed or substantially modified flags/protocols.
   - `SECURITY_RELEVANT`: Changes to authentication or sandboxing behavior.
3. **Capture Fixtures**: Update `tests/fixtures/agy/` with sanitized samples of the new version's output.
4. **Update Adapters**: Adjust `internal/upstream/agy/` to bridge differences without altering core AGY++ logic.
5. **Document**: Record results in `docs/upstream/UPSTREAM.md` and `config/upstream.yaml`.
