package agy

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// TurnResult is returned by RunTurn after each prompt turn.
type TurnResult struct {
	ConversationID string
	Response       string
	Status         string
	Error          string
}

// RunTurn executes a single prompt turn via agy --output-format stream-json,
// streaming text_delta fragments to the writer in real-time.
// It returns the conversation_id so callers can resume with --conversation.
func (p *ProcessAdapter) RunTurn(ctx context.Context, opts TurnOptions) (*TurnResult, error) {
	args := p.buildTurnArgs(opts)

	cmd := exec.CommandContext(ctx, p.Installation.BinaryPath, args...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start agy: %w", err)
	}

	result := &TurnResult{}
	scanner := bufio.NewScanner(stdout)
	out := opts.Out
	if out == nil {
		out = os.Stdout
	}

	for scanner.Scan() {
		line := scanner.Bytes()
		evt, err := ParseStreamLine(line)
		if err != nil {
			continue // skip malformed lines
		}

		switch evt.Event {
		case "init":
			if evt.Init != nil {
				result.ConversationID = evt.Init.ConversationID
			}
		case "step_update":
			if evt.StepUpdate != nil && evt.StepUpdate.StepType == "agent_response" {
				delta := evt.StepUpdate.TextDelta
				if delta != "" {
					fmt.Fprint(out, delta)
				}
			}
		case "result":
			if evt.Result != nil {
				result.ConversationID = evt.Result.ConversationID
				result.Response = evt.Result.Response
				result.Status = evt.Result.Status
				result.Error = evt.Result.Error
			}
		}
	}

	_ = cmd.Wait()
	fmt.Fprintln(out) // trailing newline after streamed response
	return result, nil
}

// TurnOptions controls a single stream-json turn.
type TurnOptions struct {
	Prompt         string
	ConversationID string // if set, resumes via --conversation
	Continue       bool   // if set, resumes most recent via --continue
	Yolo           bool
	PlanMode       bool
	Sandbox        bool
	Model          string
	Project        string
	Effort         string
	Out            io.Writer // where to stream text_delta fragments
}

func (p *ProcessAdapter) buildTurnArgs(opts TurnOptions) []string {
	var args []string

	// YOLO / mode flags must match what was used to start the session,
	// since --dangerously-skip-permissions applies at session level.
	if opts.Yolo && p.Installation.Capabilities.SupportsYolo {
		args = append(args, "--dangerously-skip-permissions")
	}
	if opts.PlanMode && p.Installation.Capabilities.SupportsPlan {
		args = append(args, "--mode", "plan")
	}
	if opts.Sandbox && p.Installation.Capabilities.SupportsSandbox {
		args = append(args, "--sandbox")
	}
	if opts.Model != "" {
		args = append(args, "--model", opts.Model)
	}
	if opts.Project != "" {
		args = append(args, "--project", opts.Project)
	}
	if opts.Effort != "" {
		args = append(args, "--effort", opts.Effort)
	}

	// Session continuity
	if opts.ConversationID != "" {
		args = append(args, "--conversation", opts.ConversationID)
	} else if opts.Continue {
		args = append(args, "--continue")
	}

	// Prompt
	if opts.Prompt != "" {
		args = append(args, "--print", opts.Prompt)
	}

	// Stream JSON output
	args = append(args, "--output-format", "stream-json")

	return args
}

// IsContinueArg checks if an arg slice references conversation resumption.
func IsContinueArg(args []string) bool {
	for _, a := range args {
		if a == "--continue" || a == "-c" || strings.HasPrefix(a, "--conversation") {
			return true
		}
	}
	return false
}
