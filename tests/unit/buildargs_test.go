package unit

import (
	"testing"

	"github.com/cdrivex4/agy-plus-plus/internal/cli"
	"github.com/cdrivex4/agy-plus-plus/internal/policy"
	"github.com/cdrivex4/agy-plus-plus/internal/upstream/agy"
)

// makeInst creates a mock agy Installation with all capabilities enabled.
func makeInst() *agy.Installation {
	return &agy.Installation{
		BinaryPath: "agy.exe",
		Version:    "1.2.2",
		Capabilities: &agy.Capabilities{
			SupportsYolo:              true,
			SupportsPlan:              true,
			SupportsSandbox:           true,
			SupportsInteractivePrompt: true,
			SupportsStreamJSON:        true,
		},
	}
}

func containsArg(args []string, target string) bool {
	for _, a := range args {
		if a == target {
			return true
		}
	}
	return false
}

func containsSequence(args []string, key, val string) bool {
	for i, a := range args {
		if a == key && i+1 < len(args) && args[i+1] == val {
			return true
		}
	}
	return false
}

func TestBuildAgyArgs_Standard(t *testing.T) {
	flags := &cli.ParsedFlags{}
	pol := policy.NewManager(false)
	inst := makeInst()

	args := cli.BuildAgyArgs(flags, pol, inst)

	// Standard mode: no special flags
	if containsArg(args, "--dangerously-skip-permissions") {
		t.Error("standard mode must NOT include --dangerously-skip-permissions")
	}
	if containsArg(args, "--mode") {
		t.Error("standard mode must NOT include --mode")
	}
}

func TestBuildAgyArgs_Yolo(t *testing.T) {
	flags := &cli.ParsedFlags{YoloMode: true}
	pol := policy.NewManager(false)
	_ = pol.EnableYolo(policy.SourceCLIArg)
	inst := makeInst()

	args := cli.BuildAgyArgs(flags, pol, inst)

	if !containsArg(args, "--dangerously-skip-permissions") {
		t.Error("YOLO mode must include --dangerously-skip-permissions")
	}
	if containsArg(args, "--mode") {
		t.Error("YOLO mode must NOT include --mode plan")
	}
}

func TestBuildAgyArgs_PlanMode(t *testing.T) {
	flags := &cli.ParsedFlags{PlanMode: true}
	pol := policy.NewManager(false)
	pol.SetPlanMode(true)
	inst := makeInst()

	args := cli.BuildAgyArgs(flags, pol, inst)

	if containsArg(args, "--dangerously-skip-permissions") {
		t.Error("plan mode must NOT include --dangerously-skip-permissions")
	}
	if !containsSequence(args, "--mode", "plan") {
		t.Error("plan mode must include --mode plan")
	}
}

func TestBuildAgyArgs_Sandbox(t *testing.T) {
	flags := &cli.ParsedFlags{Sandbox: true}
	pol := policy.NewManager(true)
	inst := makeInst()

	args := cli.BuildAgyArgs(flags, pol, inst)

	if !containsArg(args, "--sandbox") {
		t.Error("sandbox flag must include --sandbox")
	}
}

func TestBuildAgyArgs_PrintMode(t *testing.T) {
	flags := &cli.ParsedFlags{PrintPrompt: "hello world"}
	pol := policy.NewManager(false)
	inst := makeInst()

	args := cli.BuildAgyArgs(flags, pol, inst)

	if !containsSequence(args, "--print", "hello world") {
		t.Error("print mode must include --print <prompt>")
	}
}

func TestBuildAgyArgs_PrintModeYolo(t *testing.T) {
	flags := &cli.ParsedFlags{PrintPrompt: "hello world", YoloMode: true}
	pol := policy.NewManager(false)
	_ = pol.EnableYolo(policy.SourceCLIArg)
	inst := makeInst()

	args := cli.BuildAgyArgs(flags, pol, inst)

	if !containsArg(args, "--dangerously-skip-permissions") {
		t.Error("print+yolo mode must include --dangerously-skip-permissions")
	}
	if !containsSequence(args, "--print", "hello world") {
		t.Error("print mode must include --print")
	}
}

func TestBuildAgyArgs_Continue(t *testing.T) {
	flags := &cli.ParsedFlags{Continue: true}
	pol := policy.NewManager(false)
	inst := makeInst()

	args := cli.BuildAgyArgs(flags, pol, inst)

	if !containsArg(args, "--continue") {
		t.Error("continue flag must include --continue")
	}
}

func TestBuildAgyArgs_ConversationID(t *testing.T) {
	flags := &cli.ParsedFlags{ConversationID: "abc123"}
	pol := policy.NewManager(false)
	inst := makeInst()

	args := cli.BuildAgyArgs(flags, pol, inst)

	if !containsSequence(args, "--conversation", "abc123") {
		t.Error("conversation flag must include --conversation <id>")
	}
}

func TestBuildAgyArgs_Model(t *testing.T) {
	flags := &cli.ParsedFlags{Model: "gemini-3.8-flash-high"}
	pol := policy.NewManager(false)
	inst := makeInst()

	args := cli.BuildAgyArgs(flags, pol, inst)

	if !containsSequence(args, "--model", "gemini-3.8-flash-high") {
		t.Error("model flag must include --model")
	}
}

func TestBuildAgyArgs_CapabilityMissing_Yolo(t *testing.T) {
	// AGY binary without YOLO capability — flag must NOT be injected
	flags := &cli.ParsedFlags{YoloMode: true}
	pol := policy.NewManager(false)
	_ = pol.EnableYolo(policy.SourceCLIArg)
	inst := &agy.Installation{
		BinaryPath: "agy.exe",
		Version:    "0.5.0",
		Capabilities: &agy.Capabilities{
			SupportsYolo: false, // explicitly missing
		},
	}

	args := cli.BuildAgyArgs(flags, pol, inst)

	if containsArg(args, "--dangerously-skip-permissions") {
		t.Error("must NOT inject YOLO flag when backend capability is missing")
	}
}

func TestBuildAgyArgs_InteractivePrompt(t *testing.T) {
	flags := &cli.ParsedFlags{InteractivePrompt: "review my code"}
	pol := policy.NewManager(false)
	inst := makeInst()

	args := cli.BuildAgyArgs(flags, pol, inst)

	if !containsSequence(args, "--prompt-interactive", "review my code") {
		t.Error("interactive prompt must include --prompt-interactive")
	}
}

func TestBuildAgyArgs_OutputFormat_Default(t *testing.T) {
	// Default output format "text" should NOT be passed to agy (it's the agy default too)
	flags := &cli.ParsedFlags{PrintPrompt: "test", OutputFormat: "text"}
	pol := policy.NewManager(false)
	inst := makeInst()

	args := cli.BuildAgyArgs(flags, pol, inst)

	if containsArg(args, "--output-format") {
		t.Error("default 'text' output format must NOT be explicitly passed to agy")
	}
}

func TestBuildAgyArgs_OutputFormat_JSON(t *testing.T) {
	flags := &cli.ParsedFlags{PrintPrompt: "test", OutputFormat: "json"}
	pol := policy.NewManager(false)
	inst := makeInst()

	args := cli.BuildAgyArgs(flags, pol, inst)

	if !containsSequence(args, "--output-format", "json") {
		t.Error("json output format must be passed to agy explicitly")
	}
}
