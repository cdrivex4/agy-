package agy

import (
	"context"
	"io"
	"os"
	"os/exec"
)

// ProcessAdapter supervises the execution of the upstream agy process.
type ProcessAdapter struct {
	Installation *Installation
}

// NewProcessAdapter creates a process adapter from an Installation.
func NewProcessAdapter(inst *Installation) *ProcessAdapter {
	return &ProcessAdapter{Installation: inst}
}

// InteractiveOptions holds settings for interactive execution.
type InteractiveOptions struct {
	Prompt       string
	Yolo         bool
	PlanMode     bool
	Sandbox      bool
	Model        string
	Project      string
	ExtraArgs    []string
	Stdin        io.Reader
	Stdout       io.Writer
	Stderr       io.Writer
}

// RunInteractive launches upstream agy in interactive mode attached to terminal streams.
func (p *ProcessAdapter) RunInteractive(ctx context.Context, opts InteractiveOptions) error {
	args := []string{}

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
	if opts.Prompt != "" {
		args = append(args, "-i", opts.Prompt)
	}

	args = append(args, opts.ExtraArgs...)

	cmd := exec.CommandContext(ctx, p.Installation.BinaryPath, args...)
	cmd.Stdin = opts.Stdin
	if cmd.Stdin == nil {
		cmd.Stdin = os.Stdin
	}
	cmd.Stdout = opts.Stdout
	if cmd.Stdout == nil {
		cmd.Stdout = os.Stdout
	}
	cmd.Stderr = opts.Stderr
	if cmd.Stderr == nil {
		cmd.Stderr = os.Stderr
	}

	return cmd.Run()
}

// PrintOptions holds settings for non-interactive execution.
type PrintOptions struct {
	Prompt       string
	Yolo         bool
	OutputFormat string // text, json, stream-json
	ExtraArgs    []string
	Stdout       io.Writer
	Stderr       io.Writer
}

// RunPrint launches upstream agy with -p (non-interactive).
func (p *ProcessAdapter) RunPrint(ctx context.Context, opts PrintOptions) error {
	args := []string{"-p", opts.Prompt}

	if opts.Yolo && p.Installation.Capabilities.SupportsYolo {
		args = append(args, "--dangerously-skip-permissions")
	}
	if opts.OutputFormat != "" {
		args = append(args, "--output-format", opts.OutputFormat)
	}

	args = append(args, opts.ExtraArgs...)

	cmd := exec.CommandContext(ctx, p.Installation.BinaryPath, args...)
	cmd.Stdout = opts.Stdout
	if cmd.Stdout == nil {
		cmd.Stdout = os.Stdout
	}
	cmd.Stderr = opts.Stderr
	if cmd.Stderr == nil {
		cmd.Stderr = os.Stderr
	}

	return cmd.Run()
}
