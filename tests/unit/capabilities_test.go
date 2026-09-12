package unit

import (
	"testing"

	"github.com/cdrivex4/agy-plus-plus/internal/upstream/agy"
)

func TestCapabilities_AllTrue(t *testing.T) {
	help := `
  --dangerously-skip-permissions  Auto-approve all tool permission requests
  --mode                          Set the agent execution mode (accept-edits, plan)
  --output-format                 Output format (text, json, stream-json)
  --prompt-interactive            Run interactively
  --sandbox                       Run with terminal restrictions
  mcp                             Manage MCP servers
  models                          List available models
  agents                          List available agents
`
	caps := agy.ParseCapabilities(help)

	if !caps.SupportsYolo { t.Error("expected SupportsYolo=true") }
	if !caps.SupportsPlan { t.Error("expected SupportsPlan=true") }
	if !caps.SupportsStreamJSON { t.Error("expected SupportsStreamJSON=true") }
	if !caps.SupportsInteractivePrompt { t.Error("expected SupportsInteractivePrompt=true") }
	if !caps.SupportsSandbox { t.Error("expected SupportsSandbox=true") }
	if !caps.SupportsMCP { t.Error("expected SupportsMCP=true") }
	if !caps.SupportsModels { t.Error("expected SupportsModels=true") }
	if !caps.SupportsAgents { t.Error("expected SupportsAgents=true") }
}

func TestCapabilities_Minimal(t *testing.T) {
	// A minimal older binary with no advanced flags
	caps := agy.ParseCapabilities("Usage: agy\n  -p run a prompt\n")

	if caps.SupportsYolo { t.Error("expected SupportsYolo=false for minimal binary") }
	if caps.SupportsPlan { t.Error("expected SupportsPlan=false for minimal binary") }
	if caps.SupportsSandbox { t.Error("expected SupportsSandbox=false for minimal binary") }
}

func TestCapabilities_SupportsMethod(t *testing.T) {
	caps := &agy.Capabilities{
		SupportsYolo:    true,
		SupportsPlan:    false,
		SupportsSandbox: true,
	}

	if !caps.Supports("yolo") { t.Error("Supports('yolo') should be true") }
	if !caps.Supports("skip-permissions") { t.Error("Supports('skip-permissions') alias should be true") }
	if caps.Supports("plan") { t.Error("Supports('plan') should be false") }
	if !caps.Supports("sandbox") { t.Error("Supports('sandbox') should be true") }
	if caps.Supports("unknown-feature") { t.Error("Supports('unknown-feature') should be false") }
}
