package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/cdrivex4/agy-plus-plus/internal/diagnostics"
	"github.com/cdrivex4/agy-plus-plus/internal/policy"
	"github.com/cdrivex4/agy-plus-plus/internal/upstream/agy"
	"github.com/cdrivex4/agy-plus-plus/internal/version"
)

// ParsedFlags holds all CLI flags parsed before execution.
type ParsedFlags struct {
	ShowVersion      bool
	YoloMode         bool
	PlanMode         bool
	Sandbox          bool
	PrintPrompt      string
	InteractivePrompt string
	OutputFormat     string
	Continue         bool
	ConversationID   string
	Model            string
	Project          string
	Effort           string
	ExtraArgs        []string
}

// Execute runs the AGY++ CLI application.
func Execute(ctx context.Context, args []string) int {
	// Handle bare subcommands before flag parsing
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

	// Parse flags
	flags, err := ParseFlags(args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
		printUsage()
		return 1
	}

	if flags.ShowVersion {
		fmt.Println(version.String())
		return 0
	}

	// Discover upstream AGY
	inst, err := agy.Discover()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n\nEnsure Antigravity CLI (agy) is installed and in your PATH.\n", err)
		return 1
	}

	// Build and validate policy
	pol := policy.NewManager(flags.Sandbox)
	if flags.YoloMode {
		if err := pol.EnableYolo(policy.SourceCLIArg); err != nil {
			fmt.Fprintf(os.Stderr, "Policy Error: %v\n", err)
			return 1
		}
	}
	if flags.PlanMode {
		pol.SetPlanMode(true)
	}

	// Build AGY args
	agyArgs := BuildAgyArgs(flags, pol, inst)

	// Print AGY++ mode header
	printSessionHeader(pol)

	// Execute upstream agy
	adapter := agy.NewProcessAdapter(inst)
	if err := adapter.RunWithArgs(ctx, agyArgs); err != nil {
		// agy exits non-zero when user quits normally — don't treat that as an error
		return 0
	}

	return 0
}

// BuildAgyArgs constructs the argument slice to pass directly to agy.exe.
func BuildAgyArgs(flags *ParsedFlags, pol *policy.Manager, inst *agy.Installation) []string {
	var args []string

	// Permission mode flags
	if pol.IsYoloEnabled() && inst.Capabilities.SupportsYolo {
		args = append(args, "--dangerously-skip-permissions")
	}
	if pol.IsPlanEnabled() && inst.Capabilities.SupportsPlan {
		args = append(args, "--mode", "plan")
	}
	if pol.SandboxActive() && inst.Capabilities.SupportsSandbox {
		args = append(args, "--sandbox")
	}

	// Session continuity
	if flags.Continue {
		args = append(args, "--continue")
	}
	if flags.ConversationID != "" {
		args = append(args, "--conversation", flags.ConversationID)
	}

	// Model and project selection
	if flags.Model != "" {
		args = append(args, "--model", flags.Model)
	}
	if flags.Project != "" {
		args = append(args, "--project", flags.Project)
	}
	if flags.Effort != "" {
		args = append(args, "--effort", flags.Effort)
	}

	// Print (non-interactive) mode
	if flags.PrintPrompt != "" {
		args = append(args, "--print", flags.PrintPrompt)
		if flags.OutputFormat != "" && flags.OutputFormat != "text" {
			args = append(args, "--output-format", flags.OutputFormat)
		}
		args = append(args, flags.ExtraArgs...)
		return args
	}

	// Interactive mode with optional initial prompt
	if flags.InteractivePrompt != "" {
		args = append(args, "--prompt-interactive", flags.InteractivePrompt)
	}

	args = append(args, flags.ExtraArgs...)
	return args
}

func printSessionHeader(pol *policy.Manager) {
	if pol.IsYoloEnabled() {
		fmt.Fprintln(os.Stderr, "╔══════════════════════════════════════╗")
		fmt.Fprintln(os.Stderr, "║  ⚡ AGY++  ·  YOLO MODE ACTIVE       ║")
		fmt.Fprintln(os.Stderr, "║  All tool actions auto-approved.      ║")
		fmt.Fprintln(os.Stderr, "║  Press Ctrl+C or type /quit to exit. ║")
		fmt.Fprintln(os.Stderr, "╚══════════════════════════════════════╝")
	} else if pol.IsPlanEnabled() {
		fmt.Fprintln(os.Stderr, "╔══════════════════════════════════════╗")
		fmt.Fprintln(os.Stderr, "║  📋 AGY++  ·  PLAN MODE ACTIVE       ║")
		fmt.Fprintln(os.Stderr, "║  Read-only guardrails enabled.        ║")
		fmt.Fprintln(os.Stderr, "╚══════════════════════════════════════╝")
	}
}

func printUsage() {
	fmt.Println(`AGY++ — Antigravity CLI on steroids

USAGE:
  agy++ [flags]                         Interactive session
  agy++ -y [flags]                      Interactive session in YOLO mode
  agy++ -p "prompt" [flags]             Non-interactive print mode
  agy++ [command]                       Run a companion command

FLAGS (agy++ specific):
  -y, --yolo                Enable YOLO mode (auto-approve all tool permissions)
      --plan                Enable Plan mode (read-only, no file edits)
      --sandbox             Enable terminal sandbox restrictions

FLAGS (passed directly to AGY):
  -i  "prompt"              Open interactive session with an initial prompt
  -p  "prompt"              Run a single prompt and print (non-interactive)
  -c, --continue            Resume the most recent conversation
      --conversation <id>   Resume a conversation by ID
      --model <model>       Select a specific model
      --project <project>   Select a project
      --effort <level>      Reasoning effort: low | medium | high
      --output-format       text | json | stream-json
  -v, --version             Show version and compatibility baseline
  -h, --help                Show this help message

COMPANION COMMANDS:
  agy++ diagnostics         Inspect system, AGY installation, and capabilities
  agy++ version             Show detailed version information
  agy++ help                Show this message

EXAMPLES:
  agy++                           Start a normal interactive AGY session
  agy++ -y                        Start session with YOLO auto-approval active
  agy++ -y -c                     Resume last session with YOLO active
  agy++ -p "Summarize this repo"  Non-interactive single prompt
  agy++ --plan -i "Review code"   Open in read-only plan mode
  agy++ diagnostics               Check system compatibility`)
}
