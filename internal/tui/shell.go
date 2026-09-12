// Package tui contains the terminal UI layer for AGY++.
//
// NOTE: AGY++ uses a native passthrough architecture - the upstream agy.exe
// owns the interactive terminal (TUI, raw mode, keystrokes, rendering).
// This package contains companion utilities that run before/after agy sessions
// (banners, diagnostics overlays, YOLO mode indicators).
//
// A full custom TUI (Ctrl+Y mid-session toggle) is a Phase 6 engineering task
// requiring stream-json mode integration with a Bubbletea-based terminal renderer.
package tui

import (
	"fmt"
	"io"
	"os"

	"github.com/cdrivex4/agy-plus-plus/internal/upstream/agy"
	"github.com/cdrivex4/agy-plus-plus/internal/version"
)

// PrintBanner writes the AGY++ session header to the given writer.
func PrintBanner(w io.Writer, inst *agy.Installation, yolo, plan bool) {
	fmt.Fprintf(w, "AGY++ v%s  |  Backend: Antigravity CLI v%s\n",
		version.Version, inst.Version)
	if yolo {
		fmt.Fprintln(w, "Mode: ⚡ YOLO (auto-approval active)")
	} else if plan {
		fmt.Fprintln(w, "Mode: 📋 PLAN (read-only guardrails)")
	} else {
		fmt.Fprintln(w, "Mode: STANDARD (confirmation prompts active)")
	}
	fmt.Fprintln(w, "────────────────────────────────────────")
}

// PrintYoloWarning writes a prominent YOLO warning to stderr.
func PrintYoloWarning() {
	fmt.Fprintln(os.Stderr, "╔══════════════════════════════════════════╗")
	fmt.Fprintln(os.Stderr, "║  ⚡ AGY++  ·  YOLO MODE ACTIVE           ║")
	fmt.Fprintln(os.Stderr, "║  All tool actions auto-approved.          ║")
	fmt.Fprintln(os.Stderr, "║  Start session with -y to enable YOLO.   ║")
	fmt.Fprintln(os.Stderr, "╚══════════════════════════════════════════╝")
}

// PrintPlanWarning writes a plan mode notice to stderr.
func PrintPlanWarning() {
	fmt.Fprintln(os.Stderr, "╔══════════════════════════════════════════╗")
	fmt.Fprintln(os.Stderr, "║  📋 AGY++  ·  PLAN MODE ACTIVE           ║")
	fmt.Fprintln(os.Stderr, "║  Read-only guardrails enabled.            ║")
	fmt.Fprintln(os.Stderr, "╚══════════════════════════════════════════╝")
}
