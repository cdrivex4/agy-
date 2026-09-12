package agy

import (
	"strings"
)

// Capabilities represents the detected features of an AGY binary.
type Capabilities struct {
	SupportsYolo              bool
	SupportsPlan              bool
	SupportsStreamJSON        bool
	SupportsInteractivePrompt bool
	SupportsSandbox           bool
	SupportsMCP               bool
	SupportsModels            bool
	SupportsAgents            bool
	SupportsChangelog         bool
	RawHelp                   string
}

// ParseCapabilities parses the output of `agy --help` and returns detected capabilities.
func ParseCapabilities(helpText string) *Capabilities {
	c := &Capabilities{
		RawHelp: helpText,
	}

	lower := strings.ToLower(helpText)

	if strings.Contains(helpText, "--dangerously-skip-permissions") {
		c.SupportsYolo = true
	}
	if strings.Contains(helpText, "--mode") && strings.Contains(lower, "plan") {
		c.SupportsPlan = true
	}
	if strings.Contains(lower, "stream-json") {
		c.SupportsStreamJSON = true
	}
	if strings.Contains(helpText, "--prompt-interactive") || strings.Contains(helpText, "-i") {
		c.SupportsInteractivePrompt = true
	}
	if strings.Contains(helpText, "--sandbox") {
		c.SupportsSandbox = true
	}
	if strings.Contains(helpText, "mcp") {
		c.SupportsMCP = true
	}
	if strings.Contains(helpText, "models") {
		c.SupportsModels = true
	}
	if strings.Contains(helpText, "agents") {
		c.SupportsAgents = true
	}
	if strings.Contains(helpText, "changelog") {
		c.SupportsChangelog = true
	}

	return c
}

// Supports checks if a specific feature keyword is supported.
func (c *Capabilities) Supports(feature string) bool {
	switch strings.ToLower(feature) {
	case "yolo", "skip-permissions":
		return c.SupportsYolo
	case "plan":
		return c.SupportsPlan
	case "stream-json":
		return c.SupportsStreamJSON
	case "interactive-prompt":
		return c.SupportsInteractivePrompt
	case "sandbox":
		return c.SupportsSandbox
	case "mcp":
		return c.SupportsMCP
	case "models":
		return c.SupportsModels
	case "agents":
		return c.SupportsAgents
	case "changelog":
		return c.SupportsChangelog
	default:
		return false
	}
}
