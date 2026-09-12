package commands

import (
	"fmt"
	"sort"
	"strings"

	"github.com/cdrivex4/agy-plus-plus/internal/session"
	"github.com/cdrivex4/agy-plus-plus/internal/upstream/agy"
)

// CommandContext holds references needed by command handlers during interactive sessions.
type CommandContext struct {
	Installation *agy.Installation
	Session      *session.Session
	SessionStore *session.Store
	Git          *session.GitHelper
	YoloEnabled  bool
	PlanEnabled  bool
	SetYolo      func(enable bool) error
	SetPlan      func(enable bool)
}

// Command represents an interactive slash command.
type Command struct {
	Name               string
	Aliases            []string
	Description        string
	RequiredCapability string
	Handler            func(ctx *CommandContext, args []string) (string, error)
}

// Registry stores registered commands.
type Registry struct {
	commands map[string]*Command
}

// NewRegistry initializes the command registry with core commands.
func NewRegistry() *Registry {
	r := &Registry{
		commands: make(map[string]*Command),
	}
	r.registerDefaults()
	return r
}

// Register adds a command to the registry.
func (r *Registry) Register(cmd *Command) {
	r.commands[cmd.Name] = cmd
	for _, alias := range cmd.Aliases {
		r.commands[alias] = cmd
	}
}

// Lookup finds a command by name or alias.
func (r *Registry) Lookup(name string) (*Command, bool) {
	name = strings.TrimPrefix(name, "/")
	cmd, ok := r.commands[name]
	return cmd, ok
}

// List returns a sorted list of unique commands.
func (r *Registry) List() []*Command {
	seen := make(map[string]bool)
	var list []*Command
	for _, cmd := range r.commands {
		if !seen[cmd.Name] {
			seen[cmd.Name] = true
			list = append(list, cmd)
		}
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].Name < list[j].Name
	})
	return list
}

// CheckCapability verifies that the given command is supported by the installed backend.
func (r *Registry) CheckCapability(cmd *Command, caps *agy.Capabilities) error {
	if cmd.RequiredCapability == "" {
		return nil
	}
	if caps == nil || !caps.Supports(cmd.RequiredCapability) {
		return fmt.Errorf("command '/%s' is unavailable with current AGY backend (missing capability: %s)",
			cmd.Name, cmd.RequiredCapability)
	}
	return nil
}

func (r *Registry) registerDefaults() {
	r.Register(&Command{
		Name:        "help",
		Aliases:     []string{"h", "?"},
		Description: "Show available slash commands and descriptions",
		Handler: func(ctx *CommandContext, args []string) (string, error) {
			var b strings.Builder
			b.WriteString("=== Available AGY++ Commands ===\n")
			for _, cmd := range r.List() {
				aliases := ""
				if len(cmd.Aliases) > 0 {
					aliases = fmt.Sprintf(" (aliases: /%s)", strings.Join(cmd.Aliases, ", /"))
				}
				b.WriteString(fmt.Sprintf("  /%-12s %s%s\n", cmd.Name, cmd.Description, aliases))
			}
			return b.String(), nil
		},
	})

	r.Register(&Command{
		Name:        "status",
		Description: "Display current session status, approval mode, and model",
		Handler: func(ctx *CommandContext, args []string) (string, error) {
			var b strings.Builder
			b.WriteString("=== AGY++ Session Status ===\n")
			if ctx.Session != nil {
				b.WriteString(fmt.Sprintf("Session ID:    %s\n", ctx.Session.ID))
				b.WriteString(fmt.Sprintf("Workspace:     %s\n", ctx.Session.Workspace))
				b.WriteString(fmt.Sprintf("Checkpoints:   %d saved\n", len(ctx.Session.Checkpoints)))
			}
			yoloStr := "DISABLED (Prompts for confirmation)"
			if ctx.YoloEnabled {
				yoloStr = "ENABLED (Auto-approving all actions)"
			}
			b.WriteString(fmt.Sprintf("YOLO Mode:     %s\n", yoloStr))

			planStr := "DISABLED"
			if ctx.PlanEnabled {
				planStr = "ENABLED (Read-only execution)"
			}
			b.WriteString(fmt.Sprintf("Plan Mode:     %s\n", planStr))

			if ctx.Installation != nil {
				b.WriteString(fmt.Sprintf("AGY Backend:   v%s (%s)\n", ctx.Installation.Version, ctx.Installation.BinaryPath))
			}
			return b.String(), nil
		},
	})

	r.Register(&Command{
		Name:               "yolo",
		Aliases:            []string{"y"},
		Description:        "Toggle YOLO mode (automatic tool permission approval)",
		RequiredCapability: "yolo",
		Handler: func(ctx *CommandContext, args []string) (string, error) {
			if ctx.SetYolo == nil {
				return "", fmt.Errorf("YOLO mode controller unavailable")
			}
			newVal := !ctx.YoloEnabled
			if err := ctx.SetYolo(newVal); err != nil {
				return "", err
			}
			if newVal {
				return "[!] YOLO MODE ENABLED: Tool actions will be approved automatically without interactive prompts.", nil
			}
			return "[*] YOLO MODE DISABLED: Standard permission confirmation prompts restored.", nil
		},
	})

	r.Register(&Command{
		Name:               "plan",
		Description:        "Toggle Plan mode (read-only planning guardrails)",
		RequiredCapability: "plan",
		Handler: func(ctx *CommandContext, args []string) (string, error) {
			if ctx.SetPlan == nil {
				return "", fmt.Errorf("plan mode controller unavailable")
			}
			newVal := !ctx.PlanEnabled
			ctx.SetPlan(newVal)
			if newVal {
				return "[*] PLAN MODE ENABLED: Read-only guardrails active. File edits and dangerous commands are blocked.", nil
			}
			return "[*] PLAN MODE DISABLED: Standard execution restored.", nil
		},
	})

	r.Register(&Command{
		Name:        "checkpoint",
		Aliases:     []string{"cp"},
		Description: "Save a conversation and Git state checkpoint",
		Handler: func(ctx *CommandContext, args []string) (string, error) {
			if ctx.Session == nil || ctx.SessionStore == nil {
				return "", fmt.Errorf("no active session to checkpoint")
			}
			note := strings.TrimSpace(strings.Join(args, " "))
			if note == "" {
				note = "Manual checkpoint"
			}
			var commit, diff string
			if ctx.Git != nil && ctx.Git.IsGitRepo() {
				commit, _ = ctx.Git.GetHeadCommit()
				diff, _ = ctx.Git.GetDiff()
			}
			cp, err := ctx.SessionStore.AddCheckpoint(ctx.Session, note, commit, diff)
			if err != nil {
				return "", fmt.Errorf("failed to save checkpoint: %w", err)
			}
			res := fmt.Sprintf("✓ Checkpoint saved [ID: %s]\n  Note: %s", cp.ID, cp.Note)
			if cp.GitCommit != "" {
				res += fmt.Sprintf("\n  Git Commit: %s", cp.GitCommit)
			}
			return res, nil
		},
	})

	r.Register(&Command{
		Name:        "diff",
		Description: "Show uncommitted Git diff in the current workspace",
		Handler: func(ctx *CommandContext, args []string) (string, error) {
			if ctx.Git == nil || !ctx.Git.IsGitRepo() {
				return "Workspace is not a Git repository or Git is unavailable.", nil
			}
			diff, err := ctx.Git.GetDiff()
			if err != nil {
				return "", fmt.Errorf("failed to read git diff: %w", err)
			}
			if strings.TrimSpace(diff) == "" {
				return "Workspace is clean (no uncommitted diffs).", nil
			}
			return diff, nil
		},
	})

	r.Register(&Command{
		Name:        "sessions",
		Description: "List saved AGY++ sessions",
		Handler: func(ctx *CommandContext, args []string) (string, error) {
			if ctx.SessionStore == nil {
				return "Session store unavailable.", nil
			}
			sessions, err := ctx.SessionStore.List()
			if err != nil {
				return "", fmt.Errorf("failed to list sessions: %w", err)
			}
			if len(sessions) == 0 {
				return "No saved sessions found.", nil
			}
			var b strings.Builder
			b.WriteString("=== Saved AGY++ Sessions ===\n")
			for _, s := range sessions {
				marker := " "
				if ctx.Session != nil && ctx.Session.ID == s.ID {
					marker = "*"
				}
				b.WriteString(fmt.Sprintf("%s [%s] %s (%s) - %d checkpoints\n",
					marker, s.ID, s.Title, s.UpdatedAt.Format("2006-01-02 15:04"), len(s.Checkpoints)))
			}
			return b.String(), nil
		},
	})

	r.Register(&Command{
		Name:        "compact",
		Description: "Summarize and compact current context while preserving raw transcript",
		Handler: func(ctx *CommandContext, args []string) (string, error) {
			return "[*] Context compacting requested. Raw session transcript remains safely preserved on disk.", nil
		},
	})

	r.Register(&Command{
		Name:        "quit",
		Aliases:     []string{"exit", "q"},
		Description: "Exit the AGY++ session",
		Handler: func(ctx *CommandContext, args []string) (string, error) {
			return "Exiting AGY++ session. Goodbye!", nil
		},
	})
}
