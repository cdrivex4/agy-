// Package tui provides the AGY++ interactive REPL shell.
//
// Architecture:
//   - Each turn is executed via agy --print --output-format stream-json (one process per turn)
//   - The conversation_id is captured from the init event and passed via --conversation on the next turn
//   - YOLO flag is applied per turn, so toggling Ctrl+Y takes effect on the VERY NEXT prompt
//   - AGY++ owns stdin in raw mode so Ctrl+Y is intercepted before agy sees input
package tui

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/term"

	"github.com/cdrivex4/agy-plus-plus/internal/policy"
	"github.com/cdrivex4/agy-plus-plus/internal/session"
	"github.com/cdrivex4/agy-plus-plus/internal/upstream/agy"
	"github.com/cdrivex4/agy-plus-plus/internal/version"
)

// Shell is the AGY++ interactive session.
type Shell struct {
	Installation   *agy.Installation
	Policy         *policy.Manager
	Session        *session.Session
	SessionStore   *session.Store
	Git            *session.GitHelper
	ConversationID string // tracks agy conversation for resume
	InitialOpts    InitialOpts
}

// InitialOpts are the launch-time options for the shell.
type InitialOpts struct {
	ConversationID string
	Continue       bool
	Model          string
	Project        string
	Effort         string
	InitialPrompt  string
}

// NewShell creates a new interactive shell session.
func NewShell(inst *agy.Installation, pol *policy.Manager, sess *session.Session, store *session.Store, opts InitialOpts) *Shell {
	cwd, _ := os.Getwd()
	return &Shell{
		Installation: inst,
		Policy:       pol,
		Session:      sess,
		SessionStore: store,
		Git:          session.NewGitHelper(cwd),
		InitialOpts:  opts,
	}
}

// Run starts the main REPL loop.
func (s *Shell) Run(ctx context.Context) error {
	s.printBanner()

	adapter := agy.NewProcessAdapter(s.Installation)

	// If the user resumed via --conversation or --continue, record it
	s.ConversationID = s.InitialOpts.ConversationID

	// Handle an initial prompt if provided via -i flag
	if s.InitialOpts.InitialPrompt != "" {
		if err := s.runTurn(ctx, adapter, s.InitialOpts.InitialPrompt); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		line, err := s.readInput()
		if err != nil {
			if err == io.EOF {
				fmt.Fprintln(os.Stderr, "\nGoodbye.")
				return nil
			}
			return err
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Slash commands handled by AGY++ (not sent to agy)
		if strings.HasPrefix(line, "/") {
			if quit := s.handleSlashCommand(line); quit {
				return nil
			}
			continue
		}

		// Regular prompt — send to agy via stream-json
		if err := s.runTurn(ctx, adapter, line); err != nil {
			fmt.Fprintf(os.Stderr, "\n[Error: %v]\n", err)
		}
	}
}

func (s *Shell) runTurn(ctx context.Context, adapter *agy.ProcessAdapter, prompt string) error {
	opts := agy.TurnOptions{
		Prompt:         prompt,
		ConversationID: s.ConversationID,
		Continue:       s.InitialOpts.Continue && s.ConversationID == "",
		Yolo:           s.Policy.IsYoloEnabled(),
		PlanMode:       s.Policy.IsPlanEnabled(),
		Model:          s.InitialOpts.Model,
		Project:        s.InitialOpts.Project,
		Effort:         s.InitialOpts.Effort,
		Out:            os.Stdout,
	}

	result, err := adapter.RunTurn(ctx, opts)
	if err != nil {
		return err
	}

	// Track conversation ID so all future turns resume correctly
	if result != nil && result.ConversationID != "" {
		s.ConversationID = result.ConversationID
		// Persist to session
		if s.Session != nil {
			s.Session.AgySessionID = result.ConversationID
			if s.SessionStore != nil {
				_ = s.SessionStore.Save(s.Session)
			}
		}
	}

	return nil
}

// readInput reads a line of input from stdin in raw mode so Ctrl+Y is intercepted.
func (s *Shell) readInput() (string, error) {
	s.printPrompt()

	fd := int(os.Stdin.Fd())

	// Non-interactive (piped stdin) — fall back to simple line reader
	if !term.IsTerminal(fd) {
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			return scanner.Text(), nil
		}
		return "", io.EOF
	}

	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return "", err
	}
	defer func() { _ = term.Restore(fd, oldState) }()

	var buf []byte
	readBuf := make([]byte, 16)

	for {
		n, err := os.Stdin.Read(readBuf)
		if err != nil {
			return "", err
		}

		for i := 0; i < n; i++ {
			b := readBuf[i]

			// Ctrl+C or Ctrl+D on empty line → exit
			if b == 0x03 || (b == 0x04 && len(buf) == 0) {
				_ = term.Restore(fd, oldState)
				return "", io.EOF
			}

			// Ctrl+Y (0x19) → toggle YOLO, update prompt in place
			if b == 0x19 {
				_ = term.Restore(fd, oldState)
				s.toggleYolo()
				// Reprint prompt with current buf
				s.printPrompt()
				fmt.Print(string(buf))
				oldState, _ = term.MakeRaw(fd)
				continue
			}

			// Enter
			if b == '\r' || b == '\n' {
				_ = term.Restore(fd, oldState)
				fmt.Println()
				return string(buf), nil
			}

			// Backspace
			if b == 0x08 || b == 0x7F {
				if len(buf) > 0 {
					buf = buf[:len(buf)-1]
					fmt.Print("\b \b")
				}
				continue
			}

			// Normal printable character
			if b >= 0x20 {
				buf = append(buf, b)
				fmt.Print(string(b))
			}
		}
	}
}

func (s *Shell) toggleYolo() {
	if s.Policy.IsYoloEnabled() {
		s.Policy.DisableYolo()
		fmt.Print("\r\n")
		fmt.Println("[*] YOLO MODE OFF — next prompt will require confirmation for tool actions.")
	} else {
		_ = s.Policy.EnableYolo(policy.SourceKeyboardShortcut)
		fmt.Print("\r\n")
		fmt.Println("[!] YOLO MODE ON  — next prompt auto-approves all tool actions.")
	}
}

func (s *Shell) printBanner() {
	fmt.Fprintf(os.Stderr, "AGY++ v%s  ·  Antigravity CLI v%s\n", version.Version, s.Installation.Version)
	fmt.Fprintln(os.Stderr, "Type /help for commands  |  Ctrl+Y to toggle YOLO  |  /quit to exit")
	fmt.Fprintln(os.Stderr, strings.Repeat("─", 56))
}

func (s *Shell) printPrompt() {
	mode := "[STD] "
	if s.Policy.IsYoloEnabled() {
		mode = "[⚡ YOLO] "
	} else if s.Policy.IsPlanEnabled() {
		mode = "[PLAN] "
	}
	fmt.Printf("\n%s> ", mode)
}

// PrintBanner writes the AGY++ session header to the given writer.
func PrintBanner(w io.Writer, inst *agy.Installation, yolo, plan bool) {
	fmt.Fprintf(w, "AGY++ v%s  |  Backend: Antigravity CLI v%s\n", version.Version, inst.Version)
	if yolo {
		fmt.Fprintln(w, "Mode: ⚡ YOLO (auto-approval active)")
	} else if plan {
		fmt.Fprintln(w, "Mode: 📋 PLAN (read-only guardrails)")
	} else {
		fmt.Fprintln(w, "Mode: STANDARD (confirmation prompts active)")
	}
	fmt.Fprintln(w, strings.Repeat("─", 42))
}

// PrintYoloWarning writes a prominent YOLO warning to stderr.
func PrintYoloWarning() {
	fmt.Fprintln(os.Stderr, "╔══════════════════════════════════════════╗")
	fmt.Fprintln(os.Stderr, "║  ⚡ AGY++  ·  YOLO MODE ACTIVE           ║")
	fmt.Fprintln(os.Stderr, "║  All tool actions auto-approved.          ║")
	fmt.Fprintln(os.Stderr, "╚══════════════════════════════════════════╝")
}

// PrintPlanWarning writes a plan mode notice to stderr.
func PrintPlanWarning() {
	fmt.Fprintln(os.Stderr, "╔══════════════════════════════════════════╗")
	fmt.Fprintln(os.Stderr, "║  📋 AGY++  ·  PLAN MODE ACTIVE           ║")
	fmt.Fprintln(os.Stderr, "║  Read-only guardrails enabled.            ║")
	fmt.Fprintln(os.Stderr, "╚══════════════════════════════════════════╝")
}

func (s *Shell) handleSlashCommand(line string) (quit bool) {
	parts := strings.Fields(line)
	cmd := strings.TrimPrefix(parts[0], "/")
	args := parts[1:]

	switch cmd {
	case "help", "h", "?":
		fmt.Println(slashHelp)
	case "status":
		s.printStatus()
	case "yolo", "y":
		s.toggleYolo()
	case "plan":
		newVal := !s.Policy.IsPlanEnabled()
		s.Policy.SetPlanMode(newVal)
		if newVal {
			fmt.Println("[*] Plan mode ON  — read-only guardrails active.")
		} else {
			fmt.Println("[*] Plan mode OFF — standard execution restored.")
		}
	case "checkpoint", "cp":
		note := strings.Join(args, " ")
		if note == "" { note = "Manual checkpoint" }
		var commit, diff string
		if s.Git.IsGitRepo() {
			commit, _ = s.Git.GetHeadCommit()
			diff, _ = s.Git.GetDiff()
		}
		if s.Session != nil && s.SessionStore != nil {
			cp, err := s.SessionStore.AddCheckpoint(s.Session, note, commit, diff)
			if err != nil {
				fmt.Printf("Error saving checkpoint: %v\n", err)
			} else {
				fmt.Printf("✓ Checkpoint saved [%s] — %s\n", cp.ID, note)
			}
		} else {
			fmt.Println("No active session to checkpoint.")
		}
	case "diff":
		if !s.Git.IsGitRepo() {
			fmt.Println("Not a git repo.")
			return false
		}
		d, _ := s.Git.GetDiff()
		if strings.TrimSpace(d) == "" {
			fmt.Println("No uncommitted changes.")
		} else {
			fmt.Println(d)
		}
	case "sessions":
		if s.SessionStore == nil { fmt.Println("Session store unavailable."); return false }
		list, _ := s.SessionStore.List()
		if len(list) == 0 { fmt.Println("No saved sessions."); return false }
		for _, ss := range list {
			marker := " "
			if s.Session != nil && ss.ID == s.Session.ID { marker = "*" }
			fmt.Printf("%s [%s] %s (%s) %d checkpoints\n",
				marker, ss.ID, ss.Title, ss.UpdatedAt.Format("2006-01-02 15:04"), len(ss.Checkpoints))
		}
	case "quit", "exit", "q":
		fmt.Println("Goodbye.")
		return true
	default:
		fmt.Printf("Unknown command: /%s — type /help for available commands.\n", cmd)
	}
	return false
}

func (s *Shell) printStatus() {
	fmt.Printf("Session ID:     %s\n", func() string {
		if s.Session != nil { return s.Session.ID }
		return "(none)"
	}())
	fmt.Printf("Conv ID:        %s\n", func() string {
		if s.ConversationID != "" { return s.ConversationID }
		return "(new)"
	}())
	fmt.Printf("YOLO Mode:      %v\n", s.Policy.IsYoloEnabled())
	fmt.Printf("Plan Mode:      %v\n", s.Policy.IsPlanEnabled())
	fmt.Printf("AGY Backend:    v%s\n", s.Installation.Version)
}

const slashHelp = `AGY++ Slash Commands:
  /help, /h        Show this help
  /status          Show session info and active mode
  /yolo, /y        Toggle YOLO mode (or press Ctrl+Y anytime)
  /plan            Toggle Plan mode (read-only guardrails)
  /checkpoint, /cp Save a conversation + git state checkpoint
  /diff            Show uncommitted git diff
  /sessions        List saved AGY++ sessions
  /quit, /q        Exit the session`
