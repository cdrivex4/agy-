# ADR-0001: Language Selection

## Status
Accepted

## Context
AGY++ requires a language that produces single, statically linked native executables for Windows (`agy++.exe`), Linux, and macOS without requiring a separate runtime (like Python or Node.js) on end-user machines. It needs robust concurrent process supervision (for communicating with the upstream AGY CLI, handling streaming NDJSON, PTY control when needed), low memory footprint, fast startup time, and strong TUI support.

Options considered:
1. **Go (Golang)**: Recommended in the master engineering instruction. Native cross-compilation, strong Windows support, single static binary, lightweight goroutines for process I/O and TUI rendering, excellent ecosystem (Cobra, Bubbletea/tview).
2. **Rust**: Excellent performance and memory safety, but steeper learning curve and slower build iterations for rapid prototype and adapter development.
3. **TypeScript / Node.js**: Excellent rapid prototyping, but requires bundling large runtimes (e.g. SEA / Bun / pkg) and has higher startup latency.

## Decision
We select **Go (v1.23+)** as the primary implementation language for AGY++.

## Consequences
- Single binary distribution with zero runtime dependencies.
- Standard Go toolchain conventions (`cmd/agy-plus-plus`, `internal/...`).
- Easy integration with GitHub Actions for cross-compiling Windows (amd64/arm64), Linux, and macOS binaries.
