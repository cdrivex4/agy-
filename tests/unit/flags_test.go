package unit

import (
	"testing"

	"github.com/cdrivex4/agy-plus-plus/internal/cli"
)

func TestParseFlags_Defaults(t *testing.T) {
	flags, err := cli.ParseFlags([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if flags.YoloMode { t.Error("YoloMode should be false by default") }
	if flags.PlanMode { t.Error("PlanMode should be false by default") }
	if flags.Sandbox { t.Error("Sandbox should be false by default") }
	if flags.PrintPrompt != "" { t.Error("PrintPrompt should be empty by default") }
	if flags.Continue { t.Error("Continue should be false by default") }
	if flags.OutputFormat != "text" { t.Errorf("OutputFormat should default to 'text', got: %s", flags.OutputFormat) }
}

func TestParseFlags_YoloLong(t *testing.T) {
	flags, _ := cli.ParseFlags([]string{"--yolo"})
	if !flags.YoloMode {
		t.Error("expected YoloMode=true from --yolo")
	}
}

func TestParseFlags_YoloShort(t *testing.T) {
	flags, _ := cli.ParseFlags([]string{"-y"})
	if !flags.YoloMode {
		t.Error("expected YoloMode=true from -y")
	}
}

func TestParseFlags_PlanMode(t *testing.T) {
	flags, _ := cli.ParseFlags([]string{"--plan"})
	if !flags.PlanMode {
		t.Error("expected PlanMode=true from --plan")
	}
}

func TestParseFlags_PrintShort(t *testing.T) {
	flags, _ := cli.ParseFlags([]string{"-p", "hello world"})
	if flags.PrintPrompt != "hello world" {
		t.Errorf("expected PrintPrompt='hello world', got: %q", flags.PrintPrompt)
	}
}

func TestParseFlags_PrintLong(t *testing.T) {
	flags, _ := cli.ParseFlags([]string{"--print", "test prompt"})
	if flags.PrintPrompt != "test prompt" {
		t.Errorf("expected PrintPrompt='test prompt', got: %q", flags.PrintPrompt)
	}
}

func TestParseFlags_Interactive(t *testing.T) {
	flags, _ := cli.ParseFlags([]string{"-i", "start prompt"})
	if flags.InteractivePrompt != "start prompt" {
		t.Errorf("expected InteractivePrompt='start prompt', got: %q", flags.InteractivePrompt)
	}
}

func TestParseFlags_Continue(t *testing.T) {
	for _, arg := range []string{"--continue", "-c"} {
		flags, _ := cli.ParseFlags([]string{arg})
		if !flags.Continue {
			t.Errorf("expected Continue=true from %s", arg)
		}
	}
}

func TestParseFlags_Model(t *testing.T) {
	flags, _ := cli.ParseFlags([]string{"--model", "gemini-3.8-flash-high"})
	if flags.Model != "gemini-3.8-flash-high" {
		t.Errorf("expected Model='gemini-3.8-flash-high', got: %q", flags.Model)
	}
}

func TestParseFlags_Effort(t *testing.T) {
	flags, _ := cli.ParseFlags([]string{"--effort", "high"})
	if flags.Effort != "high" {
		t.Errorf("expected Effort='high', got: %q", flags.Effort)
	}
}

func TestParseFlags_OutputFormat(t *testing.T) {
	flags, _ := cli.ParseFlags([]string{"--output-format", "json"})
	if flags.OutputFormat != "json" {
		t.Errorf("expected OutputFormat='json', got: %q", flags.OutputFormat)
	}
}

func TestParseFlags_ConversationID(t *testing.T) {
	flags, _ := cli.ParseFlags([]string{"--conversation", "abc123"})
	if flags.ConversationID != "abc123" {
		t.Errorf("expected ConversationID='abc123', got: %q", flags.ConversationID)
	}
}

func TestParseFlags_Sandbox(t *testing.T) {
	flags, _ := cli.ParseFlags([]string{"--sandbox"})
	if !flags.Sandbox {
		t.Error("expected Sandbox=true from --sandbox")
	}
}

func TestParseFlags_YoloAndPlan(t *testing.T) {
	// Both flags should be parseable (policy manager resolves precedence at runtime)
	flags, _ := cli.ParseFlags([]string{"--yolo", "--plan"})
	if !flags.YoloMode {
		t.Error("expected YoloMode=true")
	}
	if !flags.PlanMode {
		t.Error("expected PlanMode=true")
	}
}
