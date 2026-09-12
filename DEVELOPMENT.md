# Development Guide

## Prerequisites
- Go 1.23 or newer
- Git
- Antigravity CLI (`agy.exe`) installed on your system (v1.2.2+ recommended)

## Building Locally
To compile AGY++ for your current platform:

```powershell
# Windows
go build -o agy++.exe ./cmd/agy-plus-plus

# Or using the build script
./scripts/build/build.ps1
```

```bash
# Linux / macOS
go build -o agy++ ./cmd/agy-plus-plus
./scripts/build/build.sh
```

## Running Tests
```powershell
# Run unit tests
go test ./tests/unit/...

# Run integration tests
go test ./tests/integration/...

# Run compatibility tests against fixtures
go test ./tests/compatibility/...
```

## AI Agent Workflow Loop
When contributing or extending AGY++, follow the strict engineering workflow:
1. **READ & UNDERSTAND**: Inspect documentation and requirements first.
2. **INSPECT UPSTREAM**: Test behavior directly against upstream AGY or update fixtures.
3. **PLAN & IMPLEMENT**: Create small, focused changes.
4. **TEST & VALIDATE**: Run unit, integration, and security checks.
5. **COMMIT**: Use conventional commits (`feat:`, `fix:`, `test:`, `docs:`).
