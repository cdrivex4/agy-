package tui

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/term"

	"github.com/cdrivex4/agy-plus-plus/internal/commands"
	"github.com/cdrivex4/agy-plus-plus/internal/policy"
	"github.com/cdrivex4/agy-plus-plus/internal/session"
	"github.com/cdrivex4/agy-plus-plus/internal/upstream/agy"
	"github.com/cdrivex4/agy-plus-plus/internal/version"
)

// Shell represents the interactive AGY++ terminal session.
type Shell struct {
	Installation *agy.Installation
	Policy       *policy.Manager
	Registry     *commands.Registry
	Session      *session.Session
	SessionStore *session.Store
	Git          *session.GitHelper
	Adapter      *agy.ProcessAdapter
}

// NewShell creates a new interactive shell.
func NewShell(inst *agy.Installation, pol *policy.Manager, sess *session.Session, store *session.Store) *Shell {
	cwd, _ := os.Getwd()
	git := session.NewGitHelper(cwd)
	adapter := agy.NewProcessAdapter(inst)
	reg := commands.NewRegistry()

	return &Shell{
		Installation: inst,
		Policy:       pol,
		Registry:     reg,
		Session:      sess,
		SessionStore: store,
		Git:          git,
		Adapter:      adapter,
	}
}

// Run starts the interactive REPL loop.
func (s *Shell) Run(ctx context.Context, initialPrompt string) error {
	s.printBanner()

	if initialPrompt != "" {
		fmt.Printf("\nExecuting initial prompt: %s\n", initialPrompt)
		s.executePrompt(ctx, initialPrompt)
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		line, err := s.readLineWithCtrlY()
		if err != nil {
			if err == io.EOF {
				fmt.Println("\nExiting AGY++ session.")
				return nil
			}
			return err
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Handle Slash Commands
		if strings.HasPrefix(line, "/") {
			if s.handleSlashCommand(line) {
				return nil // Quit requested
			}
			continue
		}

		// Handle Regular AI Prompts
		s.executePrompt(ctx, line)
	}
}

func (s *Shell) printBanner() {
	fmt.Println("================================================================")
	fmt.Printf("   AGY++ v%s — Antigravity CLI on Steroids\n", version.Version)
	fmt.Println("================================================================")
	fmt.Printf("  • Backend:  Google Antigravity CLI v%s\n", s.Installation.Version)
	if s.Session != nil {
		fmt.Printf("  • Session:  %s\n", s.Session.ID)
	}
	fmt.Println("  • Commands: Type /help for slash commands | /quit to exit")
	fmt.Println("  • Shortcuts: Press Ctrl+Y anytime to toggle YOLO mode")
	fmt.Println("----------------------------------------------------------------")
	s.printStatusLine()
}

func (s *Shell) printStatusLine() {
	mode := "\033[1;32m[STANDARD MODE]\033[0m"
	if s.Policy.IsYoloEnabled() {
		mode = "\033[1;33m[⚡ YOLO MODE ACTIVE]\033[0m"
	}
	if s.Policy.IsPlanEnabled() {
		mode = "\033[1;36m[📋 PLAN MODE ACTIVE]\033[0m"
	}

	fmt.Printf("\n┌─ %s\n└──> ", mode)
}

// readLineWithCtrlY reads terminal input in raw mode to capture Ctrl+Y keystrokes live.
func (s *Shell) readLineWithCtrlY() (string, error) {
	s.printStatusLine()

	fd := int(os.Stdin.Fd())
	if !term.IsTerminal(fd) {
		// Fallback for non-interactive piped stdin
		var input string
		_, err := fmt.Scanln(&input)
		return input, err
	}

	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = term.Restore(fd, oldState)
	}()

	var buf []byte
	readBuf := make([]byte, 16)

	for {
		n, err := os.Stdin.Read(readBuf)
		if err != nil {
			return "", err
		}

		for i := 0; i < n; i++ {
			b := readBuf[i]

			// Ctrl+C (0x03) or Ctrl+D (0x04)
			if b == 0x03 || (b == 0x04 && len(buf) == 0) {
				_ = term.Restore(fd, oldState)
				return "", io.EOF
			}

			// Ctrl+Y (0x19 / decimal 25) — LIVE YOLO TOGGLE
			if b == 0x19 {
				_ = term.Restore(fd, oldState)
				if s.Policy.IsYoloEnabled() {
					s.Policy.DisableYolo()
					fmt.Print("\r\n\033[1;32m[*] YOLO MODE DISABLED: Standard permission confirmation prompts restored.\033[0m\r\n")
				} else {
					_ = s.Policy.EnableYolo(policy.SourceKeyboardShortcut)
					fmt.Print("\r\n\033[1;33m[!] YOLO MODE ENABLED: Tool actions auto-approved without prompts.\033[0m\r\n")
				}
				s.printStatusLine()
				fmt.Print(string(buf))
				oldState, _ = term.MakeRaw(fd)
				continue
			}

			// Enter (\r or \n)
			if b == '\r' || b == '\n' {
				_ = term.Restore(fd, oldState)
				fmt.Println()
				return string(buf), nil
			}

			// Backspace (0x08 or 0x7F)
			if b == 0x08 || b == 0x7F {
				if len(buf) > 0 {
					buf = buf[:len(buf)-1]
					fmt.Print("\b \b")
				}
				continue
			}

			// Printable ASCII / Unicode
			if b >= 0x20 {
				buf = append(buf, b)
				fmt.Print(string(b))
			}
		}
	}
}

func (s *Shell) handleSlashCommand(line string) bool {
	parts := strings.Fields(line)
	cmdName := strings.TrimPrefix(parts[0], "/")
	args := parts[1:]

	cmd, ok := s.Registry.Lookup(cmdName)
	if !ok {
		fmt.Printf("Unknown command: /%s. Type /help for available commands.\n", cmdName)
		return false
	}

	if err := s.Registry.CheckCapability(cmd, s.Installation.Capabilities); err != nil {
		fmt.Printf("Error: %v\n", err)
		return false
	}

	ctx := &commands.CommandContext{
		Installation: s.Installation,
		Session:      s.Session,
		SessionStore: s.SessionStore,
		Git:          s.Git,
		YoloEnabled:  s.Policy.IsYoloEnabled(),
		PlanEnabled:  s.Policy.IsPlanEnabled(),
		SetYolo: func(enable bool) error {
			if enable {
				return s.Policy.EnableYolo(policy.SourceInteractiveCmd)
			}
			s.Policy.DisableYolo()
			return nil
		},
		SetPlan: func(enable bool) {
			s.Policy.SetPlanMode(enable)
		},
	}

	if cmd.Name == "quit" {
		out, _ := cmd.Handler(ctx, args)
		fmt.Println(out)
		return true // Terminate shell
	}

	out, err := cmd.Handler(ctx, args)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else if out != "" {
		fmt.Println(out)
	}

	return false
}

func (s *Shell) executePrompt(ctx context.Context, prompt string) {
	fmt.Println()
	err := s.Adapter.RunPrint(ctx, agy.PrintOptions{
		Prompt: prompt,
		Yolo:   s.Policy.IsYoloEnabled(),
	})
	if err != nil {
		fmt.Printf("\n[Error executing turn: %v]\n", err)
	}
	fmt.Println()
}
