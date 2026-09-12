package diagnostics

import (
	"fmt"
	"regexp"
	"runtime"
	"strings"

	"github.com/cdrivex4/agy-plus-plus/internal/upstream/agy"
	"github.com/cdrivex4/agy-plus-plus/internal/version"
)

var (
	// Redaction patterns for tokens, keys, and secrets
	secretPatterns = []*regexp.Regexp{
		regexp.MustCompile(`(?i)(ya29\.[0-9a-zA-Z_-]+)`),
		regexp.MustCompile(`(?i)(AIza[0-9A-Za-z-_]{35})`),
		regexp.MustCompile(`(?i)(Bearer\s+[A-Za-z0-9\-._~+/]+=*)`),
		regexp.MustCompile(`(?i)(refresh_token=)[^\s&]+`),
	}
)

// Redact sanitizes sensitive substrings from input strings.
func Redact(input string) string {
	redacted := input
	for _, pattern := range secretPatterns {
		redacted = pattern.ReplaceAllString(redacted, "[REDACTED]")
	}
	return redacted
}

// GenerateReport produces a diagnostic overview of AGY++ and the underlying AGY environment.
func GenerateReport(inst *agy.Installation) string {
	var b strings.Builder

	b.WriteString("=== AGY++ Diagnostics ===\n")
	b.WriteString(fmt.Sprintf("AGY++ Version:        %s\n", version.Version))
	b.WriteString(fmt.Sprintf("Git Commit:           %s\n", version.GitCommit))
	b.WriteString(fmt.Sprintf("Build Date:           %s\n", version.BuildDate))
	b.WriteString(fmt.Sprintf("Platform:             %s/%s\n", runtime.GOOS, runtime.GOARCH))
	b.WriteString(fmt.Sprintf("Baseline AGY Target:  %s\n\n", version.BaselineAgy))

	b.WriteString("--- Upstream AGY Installation ---\n")
	if inst == nil {
		b.WriteString("Status:               NOT FOUND (agy executable missing)\n")
		return b.String()
	}

	b.WriteString(fmt.Sprintf("Binary Path:          %s\n", inst.BinaryPath))
	b.WriteString(fmt.Sprintf("Installed Version:    %s\n", inst.Version))

	compat := "VERIFIED"
	if inst.Version != version.BaselineAgy {
		compat = "UNTESTED (Version differs from baseline)"
	}
	b.WriteString(fmt.Sprintf("Compatibility Status: %s\n\n", compat))

	b.WriteString("--- Capabilities Detected ---\n")
	if inst.Capabilities != nil {
		b.WriteString(fmt.Sprintf("YOLO (--dangerously-skip-permissions): %v\n", inst.Capabilities.SupportsYolo))
		b.WriteString(fmt.Sprintf("Plan Mode (--mode plan):               %v\n", inst.Capabilities.SupportsPlan))
		b.WriteString(fmt.Sprintf("Stream JSON (--output-format):         %v\n", inst.Capabilities.SupportsStreamJSON))
		b.WriteString(fmt.Sprintf("Interactive Prompt (-i):               %v\n", inst.Capabilities.SupportsInteractivePrompt))
		b.WriteString(fmt.Sprintf("Sandbox (--sandbox):                   %v\n", inst.Capabilities.SupportsSandbox))
		b.WriteString(fmt.Sprintf("MCP Server Management:                 %v\n", inst.Capabilities.SupportsMCP))
		b.WriteString(fmt.Sprintf("Model Discovery (models):              %v\n", inst.Capabilities.SupportsModels))
		b.WriteString(fmt.Sprintf("Agent Discovery (agents):              %v\n", inst.Capabilities.SupportsAgents))
	}

	return Redact(b.String())
}
