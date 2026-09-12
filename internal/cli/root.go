package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/cdrivex4/agy-plus-plus/internal/diagnostics"
	"github.com/cdrivex4/agy-plus-plus/internal/policy"
	"github.com/cdrivex4/agy-plus-plus/internal/session"
	"github.com/cdrivex4/agy-plus-plus/internal/tui"
	"github.com/cdrivex4/agy-plus-plus/internal/upstream/agy"
	"github.com/cdrivex4/agy-plus-plus/internal/version"
)

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

	adapter := agy.NewProcessAdapter(inst)

	// Non-interactive print mode (no shell)
	if flags.PrintPrompt != "" {
		err := adapter.RunPrint(ctx, agy.PrintOptions{
			Prompt:       flags.PrintPrompt,
			Yolo:         pol.IsYoloEnabled(),
			OutputFormat: flags.OutputFormat,
			ExtraArgs:    flags.ExtraArgs,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return 1
		}
		return 0
	}

	// Interactive AGY++ shell with live Ctrl+Y YOLO toggling
	store, _ := session.NewStore()
	cwd, _ := os.Getwd()
	sess, _ := store.CreateSession(cwd, "")

	shellOpts := tui.InitialOpts{
		ConversationID: flags.ConversationID,
		Continue:       flags.Continue,
		Model:          flags.Model,
		Project:        flags.Project,
		Effort:         flags.Effort,
		InitialPrompt:  flags.InteractivePrompt,
	}

	sh := tui.NewShell(inst, pol, sess, store, shellOpts)
	if err := sh.Run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Session error: %v\n", err)
		return 1
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

func printUsage() {
	fmt.Println(`AGY++ — Antigravity CLI on steroids

USAGE:
  agy++ [flags]                         Interactive session (with Ctrl+Y YOLO toggle)
  agy++ -y [flags]                      Start in YOLO mode
  agy++ -p "prompt" [flags]             Non-interactive single prompt
  agy++ [command]                       Run a companion command

FLAGS (agy++ specific):
  -y, --yolo                Enable YOLO mode at session start
      --plan                Enable Plan mode (read-only, no file edits)
      --sandbox             Enable terminal sandbox restrictions

FLAGS (forwarded to AGY):
  -i  "prompt"              Open session with an initial prompt
  -p  "prompt"              Run a single prompt (non-interactive)
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

IN-SESSION CONTROLS:
  Ctrl+Y                    Toggle YOLO mode live mid-session
  /yolo, /y                 Toggle YOLO mode via slash command
  /plan                     Toggle Plan mode
  /status                   Show session info and active mode
  /checkpoint, /cp          Save a conversation + Git state checkpoint
  /diff                     Show uncommitted Git diff
  /sessions                 List saved sessions
  /help                     Show all slash commands
  /quit, /q                 Exit

EXAMPLES:
  agy++                           Start interactive session (normal mode)
  agy++ -y                        Start interactive session (YOLO mode on)
  agy++ -y -c                     Resume last session in YOLO mode
  agy++ -p "Summarize this repo"  Non-interactive single prompt
  agy++ --plan -i "Review code"   Open in read-only plan mode
  agy++ diagnostics               Check system compatibility`)
}
