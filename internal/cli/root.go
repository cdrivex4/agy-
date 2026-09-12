package cli

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/cdrivex4/agy-plus-plus/internal/diagnostics"
	"github.com/cdrivex4/agy-plus-plus/internal/policy"
	"github.com/cdrivex4/agy-plus-plus/internal/upstream/agy"
	"github.com/cdrivex4/agy-plus-plus/internal/version"
)

// Execute runs the AGY++ CLI application.
func Execute(ctx context.Context, args []string) int {
	if len(args) > 1 {
		switch args[1] {
		case "help", "--help", "-h":
			printUsage()
			return 0
		case "diagnostics":
			inst, _ := agy.Discover()
			fmt.Println(diagnostics.GenerateReport(inst))
			return 0
		case "version", "--version", "-v":
			fmt.Println(version.String())
			return 0
		}
	}

	fs := flag.NewFlagSet("agy++", flag.ContinueOnError)
	showVersion := fs.Bool("version", false, "Display version information")
	yoloLong := fs.Bool("yolo", false, "Enable YOLO mode (auto-approve tool permissions)")
	yoloShort := fs.Bool("y", false, "Enable YOLO mode (short alias)")
	planMode := fs.Bool("plan", false, "Enable Plan mode (read-only planning)")
	sandbox := fs.Bool("sandbox", false, "Run with terminal restrictions enabled")
	printPrompt := fs.String("p", "", "Run prompt non-interactively and print response")
	printPromptLong := fs.String("print", "", "Run prompt non-interactively and print response")
	interactivePrompt := fs.String("i", "", "Run an initial prompt interactively")
	outputFormat := fs.String("output-format", "text", "Output format for print mode (text, json, stream-json)")

	if err := fs.Parse(args[1:]); err != nil {
		return 1
	}

	if *showVersion {
		fmt.Println(version.String())
		return 0
	}

	// 1. Discover upstream AGY
	inst, err := agy.Discover()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 1
	}

	// 2. Initialize Policy Manager
	pol := policy.NewManager(*sandbox)

	// Activate YOLO only via human CLI flag
	if *yoloLong || *yoloShort {
		if err := pol.EnableYolo(policy.SourceCLIArg); err != nil {
			fmt.Fprintf(os.Stderr, "Policy Error: %v\n", err)
			return 1
		}
		fmt.Println("[!] AGY++: YOLO MODE ENABLED (Auto-approving permissions)")
	}

	if *planMode {
		pol.SetPlanMode(true)
		fmt.Println("[*] AGY++: PLAN MODE ENABLED (Read-only execution)")
	}

	// 3. Process Execution
	adapter := agy.NewProcessAdapter(inst)

	prompt := *printPrompt
	if prompt == "" {
		prompt = *printPromptLong
	}

	if prompt != "" {
		// Non-interactive print mode
		err := adapter.RunPrint(ctx, agy.PrintOptions{
			Prompt:       prompt,
			Yolo:         pol.IsYoloEnabled(),
			OutputFormat: *outputFormat,
			ExtraArgs:    fs.Args(),
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error during execution: %v\n", err)
			return 1
		}
		return 0
	}

	// Interactive mode
	err = adapter.RunInteractive(ctx, agy.InteractiveOptions{
		Prompt:    *interactivePrompt,
		Yolo:      pol.IsYoloEnabled(),
		PlanMode:  pol.IsPlanEnabled(),
		Sandbox:   pol.SandboxActive(),
		ExtraArgs: fs.Args(),
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Session ended with error: %v\n", err)
		return 1
	}

	return 0
}

func printUsage() {
	fmt.Println(`AGY++ — Antigravity CLI on steroids

Usage:
  agy++ [flags]
  agy++ [flags] "prompt"
  agy++ [command]

Flags:
  -y, --yolo           Enable YOLO mode (auto-approve tool permissions)
      --plan           Enable Plan mode (read-only planning guardrails)
      --sandbox        Run with terminal restrictions enabled
  -i string            Run an initial prompt interactively and continue
  -p, --print string   Run prompt non-interactively and print response
      --output-format  Output format for print mode (text, json, stream-json)
  -v, --version        Display version information
  -h, --help           Show this help message

Subcommands:
  diagnostics          Display environment, AGY discovery, and capability status
  version              Show detailed version and compatibility baseline

Interactive Slash Commands (within session):
  /help                Show available in-session commands
  /status              Display session status, active model, and YOLO mode
  /yolo                Toggle YOLO mode
  /plan                Toggle plan mode
  /diagnostics         Display system diagnostics
  /quit                Exit the session`)
}
